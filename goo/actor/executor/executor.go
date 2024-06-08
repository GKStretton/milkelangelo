package executor

import (
	"log"
	"os"

	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/dark/services/goo/events"
)

var l = log.New(os.Stdout, "[executor] ", log.Flags())

type Executor interface {
	// Final execution
	Execute()
	// Get state as expected after execution
	PredictOutcome(state *machinepb.StateReport) *machinepb.StateReport
	// Called during a voting round to issue some preemptive action, e.g. move throughout the vote.
	Preempt()
	// For debug
	String() string
}

// RunExecutorNonBlocking is a utility function to return an await channel and
// the expected state upon completion of an Executor
func RunExecutorNonBlocking(e Executor) (completionCh chan struct{}, predictedState *machinepb.StateReport) {
	predictedState = e.PredictOutcome(events.GetLatestStateReportCopy())

	completionCh = make(chan struct{})
	go func() {
		l.Printf("beginning execution: %s\n", e)
		e.Execute()
		l.Printf("completed execution: %s\n", e)

		completionCh <- struct{}{}
		close(completionCh)
	}()

	return completionCh, predictedState
}
