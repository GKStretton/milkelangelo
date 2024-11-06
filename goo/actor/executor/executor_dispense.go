package executor

import (
	"fmt"
	"time"

	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/dark/services/goo/events"
	"github.com/gkstretton/dark/services/goo/types"
)

type dispenseExecutor struct {
	X float32
	Y float32
}

func NewDispenseExecutor(d *types.DispenseDecision) *dispenseExecutor {
	if d == nil {
		return nil
	}
	return &dispenseExecutor{
		X: d.X,
		Y: d.Y,
	}
}

func (e *dispenseExecutor) Preempt() {
	goTo(e.X, e.Y)
}

func (e *dispenseExecutor) PredictOutcome(state *machinepb.StateReport) *machinepb.StateReport {
	state.PipetteState.DispenseRequestNumber++
	state.MovementDetails.TargetXUnit = e.X
	state.MovementDetails.TargetYUnit = e.Y

	state.PipetteState.VolumeTargetUl -= getDispenseVolume()
	if state.PipetteState.VolumeTargetUl < 1 {
		state.PipetteState.Spent = true
	} else {
		state.PipetteState.Spent = false
	}

	// status will change too, but not used in decision making
	return state
}

func (e *dispenseExecutor) Execute() {
	l.Printf("Going to %f, %f\n", e.X, e.Y)
	goTo(e.X, e.Y)

	// reducing to 100ms to make it more snappy
	time.Sleep(time.Millisecond * 100)

	<-events.ConditionWaiter(func(sr *machinepb.StateReport) bool {
		return sr.Status == machinepb.Status_WAITING_FOR_DISPENSE
	})

	l.Println("dispensing...")
	dispenseBlocking()
	l.Println("done...")
}

func (e *dispenseExecutor) String() string {
	return fmt.Sprintf("dispenseExecutor (x: %.3f, y: %.3f)", e.X, e.Y)
}

// call dispense, and observe transition (-> dispensing -> not dispensing)
func dispenseBlocking() {
	a1 := events.ConditionWaiter(func(sr *machinepb.StateReport) bool {
		return sr.Status == machinepb.Status_DISPENSING
	})
	time.Sleep(time.Millisecond * 250)

	err := dispense()
	if err != nil {
		l.Printf("dispense error: %v\n", err)
		innerErr := dispense()
		if innerErr != nil {
			l.Printf("dispense error after retry: %v\n", innerErr)
		}
	}

	l.Println("waiting for DISPENSING")
	<-a1

	a2 := events.ConditionWaiter(func(sr *machinepb.StateReport) bool {
		return sr.Status != machinepb.Status_DISPENSING
	})
	l.Println("waiting for NOT DISPENSING")
	<-a2
}
