package decider

import (
	"log"
	"os"

	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/dark/services/goo/types"
)

var l = log.New(os.Stdout, "[decider] ", log.Flags())

type Decider interface {
	DecideCollection(predictedState *machinepb.StateReport) *types.CollectionDecision
	DecideDispense(predictedState *machinepb.StateReport) *types.DispenseDecision
}
