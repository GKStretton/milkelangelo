package control

import (
	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/dark/services/goo/events"
	"github.com/gkstretton/dark/services/goo/types"
)

// listenForEbsCommands runs commands received from the EBS. Each command runs
// in its own goroutine so the listener isn't blocked; commands that arrive
// while another is in progress are rejected with ErrBusy.
func listenForEbsCommands() {
	c := ebsApi.SubscribeMessages()
	defer ebsApi.UnsubscribeMessages(c)

	for msg := range c {
		switch msg.Type {
		case types.EbsCollectionRequest:
			if msg.CollectionRequest == nil {
				continue
			}
			vial := msg.CollectionRequest.Id
			go func() { logErr("collect", Collect(vial)) }()
		case types.EbsDispenseRequest:
			if msg.DispenseRequest == nil {
				continue
			}
			x, y := msg.DispenseRequest.X, msg.DispenseRequest.Y
			go func() { logErr("dispense", Dispense(x, y)) }()
		case types.EbsGoToRequest:
			if msg.GoToRequest == nil {
				continue
			}
			logErr("goto", GoTo(msg.GoToRequest.X, msg.GoToRequest.Y))
		}
	}
}

func logErr(action string, err error) {
	if err != nil {
		l.Printf("%s rejected: %v\n", action, err)
	}
}

// watchStateReports keeps the EBS's view of what's allowed in sync with the
// machine.
func watchStateReports() {
	c := events.Subscribe()
	defer events.Unsubscribe(c)

	for sr := range c {
		updateEbsState(sr)
	}
}

// updateEbsState tells the EBS which commands are currently accepted.
func updateEbsState(sr *machinepb.StateReport) {
	if ebsApi == nil {
		return
	}

	lock.Lock()
	ready := enabled && !busy && sr != nil && isAwake(sr)
	e := enabled
	lock.Unlock()

	spent := sr.GetPipetteState().GetSpent()

	ebsApi.UpdateState(func(state *types.GooState) {
		state.ControlEnabled = e
		state.WaitingForCollection = ready && spent
		state.WaitingForDispense = ready && !spent
	})
}
