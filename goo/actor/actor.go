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

type decision struct {
	e   executor.Executor
	err error
}

var lock *ActorLock = &ActorLock{}

var exitCh chan struct{} = make(chan struct{}, 1)

func shouldExit() bool {
	select {
	case <-exitCh:
		return true
	default:
		return false
	}
}

func Setup() {
	subscribeToBrokerTopics()
}

// LaunchActor is launched to control a session after the canvas is prepared.
// It should effect art.
func LaunchActor(twitchApi *twitchapi.TwitchApi, actorTimeout time.Duration) error {
	if lock.Get() {
		return fmt.Errorf("actor already running")
	}
	lock.Set(true)
	defer lock.Set(false)

	// clear exit flag
	_ = shouldExit()

	fmt.Printf("Launching actor\n")

	// access state reports
	c := events.Subscribe()
	defer events.Unsubscribe(c)

	// ebs, err := ebsinterface.NewExtensionSession(time.Hour * 2)
	// if err != nil {
	// 	fmt.Printf("failed to create ebs interface in LaunchActor: %v\n", err)
	// }

	// d := decider.NewTwitchDecider(ebs, twitchApi, time.Second*5, decider.NewMockDecider())
	d := decider.NewAutoDecider(actorTimeout)

	awaitDecision := decide(d, events.GetLatestStateReportCopy(), nil)
	decision := <-awaitDecision

	for {
		if shouldExit() {
			l.Printf("exit triggered")
			break
		}
		if decision.err != nil {
			l.Printf("error in decider, exiting actor: %s\n", decision.err)
			return decision.err
		}
		if decision.e == nil {
			l.Println("saw nil decision, exiting actor")
			break
		}
		awaitCompletion, predictedCompletionState := executor.RunExecutorNonBlocking(c, decision.e)
		// ebs.UpdateCurrentAction(e)
		// ebs.UpdateUpcomingAction(nil)

		// get next action while the action is being performed
		awaitDecision = decide(d, predictedCompletionState, nil)

		<-awaitCompletion // ensure last action finished
		// ebs.UpdateCurrentAction(nil)
		decision = <-awaitDecision
	}

	return nil
}

func decide(decider decider.Decider, predictedState *machinepb.StateReport, ebs *ebsinterface.ExtensionSession) chan decision {
	c := make(chan decision)

	if predictedState == nil {
		c <- decision{
			e:   nil,
			err: fmt.Errorf("predictatedState nil"),
		}
		close(c)
		return c
	}

	go func() {
		l.Printf("making next decision...\n")
		e, err := decider.DecideNextAction(predictedState)
		if *waitForUser {
			fmt.Scanln()
		}
		if err != nil {
			l.Printf("error making next decision: %v\n", err)
		} else {
			l.Printf("made next decision: %v\n", e)
		}

		// ebs.UpdateUpcomingAction(e)
		c <- decision{
			e:   e,
			err: err,
		}
		close(c)
	}()
	return c
}
