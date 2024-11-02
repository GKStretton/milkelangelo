package decider

import (
	"log"
	"os"

	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/dark/services/goo/actor/executor"
	"github.com/gkstretton/dark/services/goo/types"
)

var l = log.New(os.Stdout, "[decider] ", log.Flags())

type Decider interface {
	decideCollection(predictedState *machinepb.StateReport) *types.CollectionDecision
	decideDispense(predictedState *machinepb.StateReport) *types.DispenseDecision
	// DecideNextAction is per decider in case there's more functions in future that vary
	// per decider. like bowl spinning or fans
	DecideNextAction(predictedState *machinepb.StateReport) (executor.Executor, error)
}
