// Package control executes collect / dispense / goto commands directly on the
// machine, on behalf of a remote controller (see publicapi).
//
// Commands are only accepted while control is enabled, the machine is awake,
// and no other command is in progress. Machine-level safety checks (vial
// validity, coordinate bounds, pipette state) are enforced here, not in the
// public api, so the public-facing side can't bypass them.
//
// Collect and Dispense validate synchronously and return an error if the
// command is rejected; an accepted command then runs in the background.
package control

import (
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"sync"
	"time"

	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/dark/services/goo/events"
	"github.com/gkstretton/dark/services/goo/types"
	"github.com/gkstretton/dark/services/goo/vialprofiles"
)

var l = log.New(os.Stdout, "[control] ", log.Flags())

var (
	ErrDisabled = errors.New("control is disabled")
	ErrBusy     = errors.New("machine is busy with another command")
)

const collectionDrops = 3

var (
	lock    sync.Mutex
	enabled bool
	busy    bool

	stateUpdater types.GooStateUpdater
)

// Start wires up control to the broker. su, if non-nil, is kept informed of
// which commands are currently accepted.
func Start(su types.GooStateUpdater) {
	stateUpdater = su

	subscribeToBrokerTopics()
	go watchStateReports()
}

// SetEnabled allows or disallows control commands.
func SetEnabled(e bool) {
	lock.Lock()
	enabled = e
	lock.Unlock()

	l.Printf("control enabled: %t\n", e)
	publishStatus()
	updatePublicState(events.GetLatestStateReportCopy())
}

func IsEnabled() bool {
	lock.Lock()
	defer lock.Unlock()
	return enabled
}

// Collect starts collecting dye from the given vial.
func Collect(vial int) error {
	err := begin(func(sr *machinepb.StateReport) error {
		if !sr.GetPipetteState().GetSpent() {
			return errors.New("pipette still holds dye, dispense first")
		}
		profile := vialprofiles.GetSystemVialProfile(vial)
		if profile == nil {
			return fmt.Errorf("no vial profile at position %d", vial)
		}
		if profile.DispenseVolumeUl <= 0 {
			return fmt.Errorf("vial profile at position %d has no dispense volume", vial)
		}
		return nil
	})
	if err != nil {
		return err
	}

	go func() {
		defer end()

		volUl := collectionDrops * int(getVialDropVolume(vial))
		l.Printf("collecting %dul from vial %d\n", volUl, vial)
		collect(vial, volUl)

		// wait for collection to start
		time.Sleep(time.Second * 1)
		<-events.ConditionWaiter(func(sr *machinepb.StateReport) bool {
			return sr.GetCollectionRequest().GetCompleted()
		})
		l.Println("collection complete")
	}()
	return nil
}

// Dispense starts moving to x, y and dispensing a drop.
func Dispense(x, y float32) error {
	err := begin(func(sr *machinepb.StateReport) error {
		if sr.GetPipetteState().GetSpent() {
			return errors.New("pipette is empty, collect first")
		}
		return validatePosition(x, y)
	})
	if err != nil {
		return err
	}

	go func() {
		defer end()

		l.Printf("going to %.3f, %.3f\n", x, y)
		goTo(x, y)

		// reducing to 100ms to make it more snappy
		time.Sleep(time.Millisecond * 100)

		<-events.ConditionWaiter(func(sr *machinepb.StateReport) bool {
			return sr.GetStatus() == machinepb.Status_WAITING_FOR_DISPENSE
		})

		l.Println("dispensing...")
		dispenseBlocking()
		l.Println("dispense complete")
	}()
	return nil
}

// GoTo moves the pipette to x, y without dispensing.
// Only allowed while holding dye, i.e. when a dispense is next.
func GoTo(x, y float32) error {
	err := begin(func(sr *machinepb.StateReport) error {
		if sr.GetPipetteState().GetSpent() {
			return errors.New("pipette is empty, collect first")
		}
		return validatePosition(x, y)
	})
	if err != nil {
		return err
	}
	defer end()

	goTo(x, y)
	return nil
}

// begin checks preconditions and marks control as busy.
func begin(check func(sr *machinepb.StateReport) error) error {
	lock.Lock()
	defer lock.Unlock()

	if !enabled {
		return ErrDisabled
	}
	if busy {
		return ErrBusy
	}

	sr := events.GetLatestStateReportCopy()
	if sr == nil {
		return errors.New("no state report available")
	}
	if !isAwake(sr) {
		return fmt.Errorf("machine not ready, status: %s", sr.GetStatus())
	}
	if err := check(sr); err != nil {
		return err
	}

	busy = true
	go updatePublicState(sr)
	return nil
}

func end() {
	lock.Lock()
	busy = false
	lock.Unlock()

	updatePublicState(events.GetLatestStateReportCopy())
}

func isAwake(sr *machinepb.StateReport) bool {
	return sr.GetStatus() >= machinepb.Status_IDLE_STATIONARY
}

// validatePosition ensures x, y is within the unit circle (the bowl)
func validatePosition(x, y float32) error {
	if math.IsNaN(float64(x)) || math.IsNaN(float64(y)) {
		return errors.New("invalid position")
	}
	if math.Hypot(float64(x), float64(y)) > 1 {
		return fmt.Errorf("position (%.3f, %.3f) outside unit circle", x, y)
	}
	return nil
}
