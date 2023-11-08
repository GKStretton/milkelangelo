package decider

import (
	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/dark/services/goo/actor/executor"
)

type twitchDecider struct{}

func NewTwitchDecider() Decider {
	return &mockDecider{}
}

func (d *twitchDecider) DecideCollection(predictedState *machinepb.StateReport) executor.Executor {
	return nil
}

func (d *twitchDecider) DecideDispense(predictedState *machinepb.StateReport) executor.Executor {
	return nil
}
