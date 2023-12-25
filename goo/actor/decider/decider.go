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
	DecideNextAction(predictedState *machinepb.StateReport) executor.Executor
}
