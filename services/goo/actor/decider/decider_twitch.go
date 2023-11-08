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
	e := executor.NewDispenseExecutor(0, 0)

	type voteStruct struct {
		x float32
		y float32
	}
	voteChan := make(chan voteStruct)
	for vote := range voteChan {
		e.X = vote.x
		e.Y = vote.y
		e.Preempt()
	}

	return e
}
