package executor

import (
	"fmt"

	"github.com/gkstretton/asol-protos/go/machinepb"
)

type dispenseExecutor struct {
	X float32
	Y float32
}

func NewDispenseExecutor(x, y float32) *dispenseExecutor {
	return &dispenseExecutor{
		X: x,
		Y: y,
	}
}

func (e *dispenseExecutor) Preempt() {
	if !goToLock.TryLock() {
		fmt.Println("preemptive goTo blocked")
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

	fmt.Printf("Going to %f, %f\n", e.X, e.Y)
	goTo(e.X, e.Y)
	fmt.Println("dispensing...")
	dispenseBlocking(c)
	fmt.Println("done...")
}

func (e *dispenseExecutor) String() string {
	return fmt.Sprintf("dispenseExecutor (x: %.3f, y: %.3f)", e.X, e.Y)
}

// call dispense, and observe transition (-> dispensing -> not dispensing)
func dispenseBlocking(c chan *machinepb.StateReport) {
	a1 := conditionWaiter(c, func(sr *machinepb.StateReport) bool {
		return sr.Status == machinepb.Status_DISPENSING
	})
	dispense()
	fmt.Println("waiting for 'DISPENSING'")
	<-a1
	fmt.Println("waiting for '!DISPENSING'")
	<-conditionWaiter(c, func(sr *machinepb.StateReport) bool {
		return sr.Status != machinepb.Status_DISPENSING
	})
}
