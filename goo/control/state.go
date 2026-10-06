package control

import (
	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/dark/services/goo/events"
	"github.com/gkstretton/dark/services/goo/types"
)

// watchStateReports keeps clients' view of what's allowed in sync with the
// machine.
func watchStateReports() {
	c := events.Subscribe()
	defer events.Unsubscribe(c)

	for sr := range c {
		updatePublicState(sr)
	}
}

// updatePublicState tells clients which commands are currently accepted.
func updatePublicState(sr *machinepb.StateReport) {
	if stateUpdater == nil {
		return
	}

	lock.Lock()
	ready := enabled && !busy && sr != nil && isAwake(sr)
	e := enabled
	lock.Unlock()

	spent := sr.GetPipetteState().GetSpent()

	stateUpdater.UpdateState(func(state *types.GooState) {
		state.ControlEnabled = e
		state.WaitingForCollection = ready && spent
		state.WaitingForDispense = ready && !spent
	})
}
