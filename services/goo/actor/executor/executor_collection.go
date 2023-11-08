package executor

import (
	"github.com/gkstretton/asol-protos/go/machinepb"
)

type collectionExecutor struct {
	vialNo int
	volUl  int
}

func NewCollectionExecutor(vialNo, volUl int) *collectionExecutor {
	return &collectionExecutor{
		vialNo,
		volUl,
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
	collect(e.vialNo, float32(e.volUl))
	<-conditionWaiter(c, func(sr *machinepb.StateReport) bool {
		return sr.CollectionRequest.Completed
	})
}
