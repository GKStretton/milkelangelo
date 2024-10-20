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
	"github.com/gkstretton/dark/services/goo/events"
	"github.com/gkstretton/dark/services/goo/session"
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

func Setup(sm *session.SessionManager) {
	subscribeToBrokerTopics(sm)
}

// LaunchActor is launched to control a session after the canvas is prepared.
// It should effect art.
func LaunchActor(twitchApi *twitchapi.TwitchApi, actorTimeout time.Duration, seed int64, testing bool) error {
	if lock.Get() {
		return fmt.Errorf("actor already running")
	}
	lock.Set(true)
	defer lock.Set(false)

	// clear exit flag
	_ = shouldExit()

	fmt.Printf("Launching actor with seed: %d\n", seed)

	endTime := time.Now().Add(actorTimeout)
	d := decider.NewAutoDecider(endTime, seed, testing)

	awaitDecision := decide(d, events.GetLatestStateReportCopy())
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
		awaitCompletion, predictedCompletionState := executor.RunExecutorNonBlocking(decision.e)

		// get next action while the action is being performed
		awaitDecision = decide(d, predictedCompletionState)

		<-awaitCompletion // ensure last action finished
		decision = <-awaitDecision
	}

	return nil
}

func decide(decider decider.Decider, predictedState *machinepb.StateReport) chan decision {
	c := make(chan decision)

	if predictedState == nil {
		c <- decision{
			e:   nil,
			err: fmt.Errorf("predictatedState nil"),
		}
		close(c)
		return c
	}

	// todo: support single user twitch control. Make new decider that falls back to the autoDecider?
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

		c <- decision{
			e:   e,
			err: err,
		}
		close(c)
	}()
	return c
}
