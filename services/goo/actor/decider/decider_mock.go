package decider

import (
	"time"

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

	return executor.NewCollectionExecutor(
		2, //empty
		15*3,
	)
}

var tempLocationTracker bool

func (d *mockDecider) DecideDispense(predictedState *machinepb.StateReport) executor.Executor {
	// simulate decision delay
	time.Sleep(time.Second * time.Duration(3))

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
