package executor

import (
	"fmt"
	"time"

	"github.com/gkstretton/asol-protos/go/machinepb"
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
	if !goToLock.TryLock() {
		l.Println("preemptive goTo blocked")
		return
	}
	defer goToLock.Unlock()
	goTo(e.X, e.Y)
}

func (e *dispenseExecutor) PredictOutcome(state *machinepb.StateReport) *machinepb.StateReport {
	state.PipetteState.DispenseRequestNumber++
	state.MovementDetails.TargetXUnit = e.X
	state.MovementDetails.TargetYUnit = e.Y

	state.PipetteState.VolumeTargetUl -= getDispenseVolume()
	if state.PipetteState.VolumeTargetUl < 1 {
		state.PipetteState.Spent = true
	}

	// status will change too, but not used in decision making
	return state
}

func (e *dispenseExecutor) Execute(c chan *machinepb.StateReport) {
	goToLock.Lock()
	defer goToLock.Unlock()

	l.Printf("Going to %f, %f\n", e.X, e.Y)
	goTo(e.X, e.Y)
	l.Println("dispensing...")
	dispenseBlocking(c)
	l.Println("done...")
}

func (e *dispenseExecutor) String() string {
	return fmt.Sprintf("dispenseExecutor (x: %.3f, y: %.3f)", e.X, e.Y)
}

// call dispense, and observe transition (-> dispensing -> not dispensing)
func dispenseBlocking(c chan *machinepb.StateReport) {
	a1 := ConditionWaiter(c, func(sr *machinepb.StateReport) bool {
		return sr.Status == machinepb.Status_DISPENSING
	})
	time.Sleep(time.Millisecond * 250)
	//! sometimes this doesn't get through...
	dispense()
	l.Println("waiting for DISPENSING")
	<-a1
	l.Println("waiting for NOT DISPENSING")
	<-ConditionWaiter(c, func(sr *machinepb.StateReport) bool {
		return sr.Status != machinepb.Status_DISPENSING
	})
}
