package decider

import (
	"log"
	"os"

	"github.com/gkstretton/asol-protos/go/machinepb"
)

var l = log.New(os.Stdout, "[decider] ", log.Flags())

type Decider interface {
	DecideCollection(predictedState *machinepb.StateReport) *CollectionDecision
	DecideDispense(predictedState *machinepb.StateReport) *DispenseDecision
}

type CollectionDecision struct {
	VialNo  int
	DropsNo int
}

type DispenseDecision struct {
	X float32
	Y float32
}
