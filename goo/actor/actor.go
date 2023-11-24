package actor

import (
	"flag"
	"fmt"

	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/dark/services/goo/actor/decider"
	"github.com/gkstretton/dark/services/goo/actor/executor"
	"github.com/gkstretton/dark/services/goo/events"
	"github.com/gkstretton/dark/services/goo/twitch"
)

var waitForUser = flag.Bool("waitForUser", false, "if true, do blocking waits at certain debug points")

// LaunchActor is launched to control a session after the canvas is prepared.
// It should effect art.
func LaunchActor(twitchApi *twitch.TwitchApi) {
	fmt.Printf("Launching actor\n")

	// access state reports
	c := events.Subscribe()
	defer events.Unsubscribe(c)

	// decider := decider.NewMockDecider()
	decider := decider.NewTwitchDecider(twitchApi)

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
		fmt.Printf("making next decision...\n")
		decision := decideNextAction(decider, predictedState)
		if *waitForUser {
			fmt.Scanln()
		}
		fmt.Printf("made   next decision: %s\n", decision)
		c <- decision
		close(c)
	}()
	return c
}

func decideNextAction(decider decider.Decider, predictedState *machinepb.StateReport) executor.Executor {
	if predictedState.PipetteState.Spent {
		fmt.Println("collection is next, deciding...")
		return decider.DecideCollection(predictedState)
	}
	fmt.Println("dispense is next, deciding...")
	return decider.DecideDispense(predictedState)
}
