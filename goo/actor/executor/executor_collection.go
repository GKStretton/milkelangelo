package executor

import (
	"fmt"

	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/dark/services/goo/actor/decider"
)

type collectionExecutor struct {
	vialNo int
	volUl  int
}

func NewCollectionExecutor(d *decider.CollectionDecision) *collectionExecutor {
	if d == nil {
		return nil
	}
	vol := d.DropsNo * getVialDropVolume(d.VialNo)
	return &collectionExecutor{
		vialNo: d.VialNo,
		volUl:  vol,
	}
}

func (e *collectionExecutor) Preempt() {}

func (e *collectionExecutor) PredictOutcome(state *machinepb.StateReport) *machinepb.StateReport {
	state.CollectionRequest.Completed = true
	state.CollectionRequest.RequestNumber++
	state.CollectionRequest.VialNumber = uint64(e.vialNo)
	state.CollectionRequest.VolumeUl = float32(e.volUl)

	state.PipetteState.VolumeTargetUl = float32(e.volUl)
	state.PipetteState.VialHeld = uint32(e.vialNo)
	state.PipetteState.Spent = false

	return state
}

func (e *collectionExecutor) Execute(c chan *machinepb.StateReport) {
	collect(e.vialNo, e.volUl)
	<-conditionWaiter(c, func(sr *machinepb.StateReport) bool {
		return sr.CollectionRequest.Completed
	})
}

func (e *collectionExecutor) String() string {
	return fmt.Sprintf("collectionExecutor (vialNo: %d, volUl: %d)", e.vialNo, e.volUl)
}
