package actor

import (
	"fmt"

	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/dark/services/goo/actor/decider"
	"github.com/gkstretton/dark/services/goo/actor/executor"
	"github.com/gkstretton/dark/services/goo/events"
)

// LaunchActor is launched to control a session after the canvas is prepared.
// It should effect art.
func LaunchActor() {
	fmt.Printf("Launching actor\n")

	// access state reports
	c := events.Subscribe()
	defer events.Unsubscribe(c)

	decider := decider.NewMockDecider()
	// decider := decider.NewTwitchDecider(),

	awaitDecision := decide(decider, events.GetLatestStateReportCopy())
	decision := <-awaitDecision

	for {
		if decision == nil {
			break
		}
		awaitCompletion, predictedCompletionState := executor.RunExecutorNonBlocking(c, decision)

		// get next action while the action is being performed
		awaitDecision = decide(decider, predictedCompletionState)

		<-awaitCompletion // ensure last action finished
		decision = <-awaitDecision
	}
}

func decide(decider decider.Decider, predictedState *machinepb.StateReport) chan executor.Executor {
	c := make(chan executor.Executor)
	go func() {
		c <- decideNextAction(decider, predictedState)
		close(c)
	}()
	return c
}

func decideNextAction(decider decider.Decider, predictedState *machinepb.StateReport) executor.Executor {
	if predictedState.PipetteState.Spent {
		return decider.DecideCollection(predictedState)
	}
	return decider.DecideDispense(predictedState)
}
