package decider

import (
	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/dark/services/goo/types"
)

type mockDecider struct{}

func NewMockDecider() Decider {
	return &mockDecider{}
}

var tempCollectionTracker int

func (d *mockDecider) DecideCollection(predictedState *machinepb.StateReport) *types.CollectionDecision {
	// Request 2 collections only
	if tempCollectionTracker >= 2 {
		return nil
	}
	tempCollectionTracker++

	//empty
	vialNo := 2
	return &types.CollectionDecision{
		VialNo:  vialNo,
		DropsNo: 3,
	}
}

var tempLocationTracker bool

func (d *mockDecider) DecideDispense(predictedState *machinepb.StateReport) *types.DispenseDecision {
	tempLocationTracker = !tempLocationTracker
	multiplier := float32(1)
	if tempLocationTracker {
		multiplier = -1
	}

	return &types.DispenseDecision{
		X: 0.5 * multiplier,
		Y: 0.5 * multiplier,
	}
}
