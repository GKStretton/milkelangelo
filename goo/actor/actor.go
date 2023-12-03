package actor

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/dark/services/goo/actor/decider"
	"github.com/gkstretton/dark/services/goo/actor/executor"
	"github.com/gkstretton/dark/services/goo/ebsinterface"
	"github.com/gkstretton/dark/services/goo/events"
	"github.com/gkstretton/dark/services/goo/twitchapi"
)

var l = log.New(os.Stdout, "[actor] ", log.Flags())
var waitForUser = flag.Bool("waitForUser", false, "if true, do blocking waits at certain debug points")

// LaunchActor is launched to control a session after the canvas is prepared.
// It should effect art.
func LaunchActor(twitchApi *twitchapi.TwitchApi, votingTimeout time.Duration) {
	fmt.Printf("Launching actor\n")

	// access state reports
	c := events.Subscribe()
	defer events.Unsubscribe(c)

	ebs, err := ebsinterface.NewExtensionSession(time.Hour * 2)
	if err != nil {
		fmt.Printf("failed to create ebs interface in LaunchActor: %v\n", err)
	}

	// todo: change to a more comprehensive auto decider
	// decider := decider.NewMockDecider()
	decider := decider.NewTwitchDecider(ebs, twitchApi, votingTimeout, decider.NewMockDecider())

	awaitDecision := decide(decider, events.GetLatestStateReportCopy())
	decision := <-awaitDecision

	for {
		if decision == nil {
			l.Println("saw nil decision, exiting actor")
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
		l.Printf("making next decision...\n")
		decision := decideNextAction(decider, predictedState)
		if *waitForUser {
			fmt.Scanln()
		}
		l.Printf("made next decision: %v\n", decision)
		c <- decision
		close(c)
	}()
	return c
}

func decideNextAction(decider decider.Decider, predictedState *machinepb.StateReport) executor.Executor {
	if predictedState.Status == machinepb.Status_SLEEPING {
		l.Println("invalid state for actor, decided nil.")
		return nil
	}
	if predictedState.PipetteState.Spent {
		l.Println("collection is next, launching decider...")
		return decider.DecideCollection(predictedState)
	}
	l.Println("dispense is next, launching decider...")
	return decider.DecideDispense(predictedState)
}
