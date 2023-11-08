package executor

import (
	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/dark/services/goo/events"
)

type Executor interface {
	// Final execution
	Execute(chan *machinepb.StateReport)
	// Get state as expected after execution
	PredictOutcome(state *machinepb.StateReport) *machinepb.StateReport
	// Called during a voting round to issue some preemptive action, e.g. move throughout the vote.
	Preempt()
}

func conditionWaiter(c chan *machinepb.StateReport, cond func(*machinepb.StateReport) bool) chan *machinepb.StateReport {
	filterChan := make(chan *machinepb.StateReport)
	go func() {
		for {
			r := <-c
			if cond(r) {
				filterChan <- r
				close(filterChan)
				return
			}
		}
	}()
	return filterChan
}

// RunExecutorNonBlocking is a utility function to return an await channel and
// the expected state upon completion of an Executor
func RunExecutorNonBlocking(c chan *machinepb.StateReport, e Executor) (completionCh chan struct{}, predictedState *machinepb.StateReport) {
	predictedState = e.PredictOutcome(events.GetLatestStateReportCopy())

	completionCh = make(chan struct{})
	go func() {
		e.Execute(c)

		completionCh <- struct{}{}
		close(completionCh)
	}()

	return completionCh, predictedState
}
