package decider

import (
	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/dark/services/goo/actor/executor"
)

type Decider interface {
	DecideCollection(predictedState *machinepb.StateReport) executor.Executor
	DecideDispense(predictedState *machinepb.StateReport) executor.Executor
}
