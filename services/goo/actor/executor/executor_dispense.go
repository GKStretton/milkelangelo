package executor

import (
	"fmt"

	"github.com/gkstretton/asol-protos/go/machinepb"
)

type dispenseExecutor struct {
	x float32
	y float32
}

func NewDispenseExecutor(x, y float32) *dispenseExecutor {
	return &dispenseExecutor{
		x,
		y,
	}
}

func (e *dispenseExecutor) Preempt() {
	if !goToLock.TryLock() {
		fmt.Println("preemptive goTo blocked")
		return
	}
	defer goToLock.Unlock()
	goTo(e.x, e.y)
}

func (e *dispenseExecutor) PredictOutcome(state *machinepb.StateReport) *machinepb.StateReport {
	state.PipetteState.DispenseRequestNumber++
	state.MovementDetails.TargetXUnit = e.x
	state.MovementDetails.TargetYUnit = e.y

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

	fmt.Printf("Going to %f, %f\n", e.x, e.y)
	goTo(e.x, e.y)
	fmt.Println("dispensing...")
	dispenseBlocking(c)
	fmt.Println("done...")
}

// call dispense, and observe transition (-> dispensing -> not dispensing)
func dispenseBlocking(c chan *machinepb.StateReport) {
	a1 := conditionWaiter(c, func(sr *machinepb.StateReport) bool {
		return sr.Status == machinepb.Status_DISPENSING
	})
	dispense()
	<-a1
	<-conditionWaiter(c, func(sr *machinepb.StateReport) bool {
		return sr.Status != machinepb.Status_DISPENSING
	})
}
