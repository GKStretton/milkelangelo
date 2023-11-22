package decider

import (
	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/dark/services/goo/actor/executor"
)

type mockDecider struct{}

func NewMockDecider() Decider {
	return &mockDecider{}
}

var tempCollectionTracker int

func (d *mockDecider) DecideCollection(predictedState *machinepb.StateReport) executor.Executor {
	// Request 2 collections only
	if tempCollectionTracker >= 2 {
		return nil
	}
	tempCollectionTracker++

	vialNo := 2
	return executor.NewCollectionExecutor(
		vialNo, //empty
		int(getVialVolume(vialNo)*3),
	)
}

var tempLocationTracker bool

func (d *mockDecider) DecideDispense(predictedState *machinepb.StateReport) executor.Executor {
	tempLocationTracker = !tempLocationTracker
	multiplier := float32(1)
	if tempLocationTracker {
		multiplier = -1
	}

	return executor.NewDispenseExecutor(
		0.5*multiplier,
		0.5*multiplier,
	)
}
