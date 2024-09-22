package scheduler

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/asol-protos/go/topics_backend"
	"github.com/gkstretton/asol-protos/go/topics_firmware"
	"github.com/gkstretton/dark/services/goo/actor"
	"github.com/gkstretton/dark/services/goo/events"
	"github.com/gkstretton/dark/services/goo/mqtt"
	"github.com/gkstretton/dark/services/goo/obs"
	"github.com/gkstretton/dark/services/goo/session"
	"github.com/gkstretton/dark/services/goo/twitchapi"
)

var sl = log.New(os.Stdout, "[session scheduler] ", log.Flags())

type SessionDescriptor struct {
	// how many minutes before session should stream start?
	streamPreStartMinutes  int
	actorDurationMinutes   int
	sessionDurationMinutes int
	runActor               bool
}

var lock *AutomationLock = &AutomationLock{}

func registerHandlers(sm *session.SessionManager, twitchApi *twitchapi.TwitchApi) {
	mqtt.Subscribe("asol/debug/runStartSequence", func(topic string, payload []byte) {
		go func() {
			fmt.Println(runStartSequence(0, false))
		}()
	})
	mqtt.Subscribe("asol/debug/runEndSequence", func(topic string, payload []byte) {
		go func() {
			fmt.Println(runEndSequence())
		}()
	})
	mqtt.Subscribe(topics_backend.TOPIC_RUN_START_SEQUENCE, func(topic string, payload []byte) {
		go func() {
			fmt.Println(runStartSequence(0, true))
		}()
	})
	mqtt.Subscribe(topics_backend.TOPIC_RUN_END_SEQUENCE, func(topic string, payload []byte) {
		go func() {
			fmt.Println(runEndSequence())
		}()
	})

	mqtt.Subscribe(topics_backend.TOPIC_RUN_FULL_SESSION, func(topic string, payload []byte) {
		go func() {
			err := RunSession(
				&SessionDescriptor{
					streamPreStartMinutes:  0,
					actorDurationMinutes:   12,
					sessionDurationMinutes: 45,
					runActor:               true,
				},
				sm, twitchApi,
			)
			if err != nil {
				fmt.Println(err)
			}
		}()
	})

	mqtt.Subscribe(topics_backend.TOPIC_RUN_MANUAL_SESSION, func(topic string, payload []byte) {
		go func() {
			err := RunSession(
				&SessionDescriptor{
					streamPreStartMinutes:  0,
					actorDurationMinutes:   12,
					sessionDurationMinutes: 45,
					runActor:               false,
				},
				sm, twitchApi,
			)
			if err != nil {
				fmt.Println(err)
			}
		}()
	})

	mqtt.Subscribe(topics_backend.TOPIC_RUN_TEST_SESSION, func(topic string, payload []byte) {
		go func() {
			n, err := strconv.Atoi(string(payload))
			if err != nil {
				fmt.Printf("couldn't get minutes from payload: %v, defaulting to 5 minutes\n", err)
				n = 5
			}
			err = RunTestSession(
				sm,
				time.Duration(n)*time.Minute,
			)
			if err != nil {
				fmt.Println(err)
			}
		}()
	})
}

func RunTestSession(sm *session.SessionManager, d time.Duration) error {
	if lock.Get() {
		return fmt.Errorf("automation already running")
	}
	lock.Set(true)
	defer lock.Set(false)

	s, err := sm.BeginSession(false)
	if err != nil {
		return err
	}

	sl.Printf("running test session %d\n", s.Id)

	sl.Println("requesting wake")
	mqtt.Publish(topics_firmware.TOPIC_WAKE, "")

	time.Sleep(time.Second * 15)
	// if not awake, abort and error
	sl.Println("checking for awake status")
	sr := events.GetLatestStateReportCopy()
	if sr.Status != machinepb.Status_IDLE_MOVING && sr.Status != machinepb.Status_IDLE_STATIONARY {
		sl.Println("invalid status, erroring")
		return fmt.Errorf("status was %s after waking but expected idle", sr.Status)
	}

	sl.Println("valid status")

	// set seed
	seed := rand.Int63()
	err = sm.SetCurrentSessionSeed(seed)
	if err != nil {
		sl.Printf("failed to set seed: %v\n", err)
	}

	sl.Println("launching actor")
	err = actor.LaunchActor(nil, d, seed, true)
	if err != nil {
		sl.Printf("actor error: %v\n", err)
	}

	sl.Println("shutting down")
	mqtt.Publish(topics_firmware.TOPIC_SHUTDOWN, "")

	waitForSleeping()

	sm.EndSession()

	return nil
}

func RunSession(
	d *SessionDescriptor,
	sm *session.SessionManager,
	twitchApi *twitchapi.TwitchApi,
) error {
	if lock.Get() {
		return fmt.Errorf("automation already running")
	}
	lock.Set(true)
	defer lock.Set(false)

	sl.Printf("running session with descriptor: %+v\n", d)

	beginTime := time.Now()

	err := runStartSequence(d.streamPreStartMinutes, true)
	if err != nil {
		return err
	}

	endTime := beginTime.Add(time.Duration(d.sessionDurationMinutes+d.streamPreStartMinutes) * time.Minute)
	drainOffsetMins := 3

	obs.SetCountdown(
		"Time until drainage:",
		endTime.Add(-time.Minute*time.Duration(drainOffsetMins)),
	)

	if d.runActor {
		seed := rand.Int63()
		err = sm.SetCurrentSessionSeed(seed)
		if err != nil {
			sl.Printf("failed to set seed: %v\n", err)
		}

		sl.Println("launching actor")
		err = actor.LaunchActor(twitchApi, time.Duration(d.actorDurationMinutes)*time.Minute, seed, false)
		if err != nil {
			sl.Println("actor error, erroring")
			mqtt.Publish(topics_backend.TOPIC_SESSION_PAUSE, "")
			// email for help
			errWrap := fmt.Errorf("actor returned error, unknown situation: %s", err)
			requestSessionIntervention(errWrap)
			return err
		}
		sl.Println("actor success")
		mqtt.Publish(topics_firmware.TOPIC_GOTO_RING_IDLE_POS, "")
	} else {
		sl.Println("ready for manual control...")
	}

	waitForTOffset(endTime, -drainOffsetMins, 0)

	err = runEndSequence()
	if err != nil {
		return err
	}

	return nil
}

func runStartSequence(streamPreStartMinutes int, realSession bool) error {
	sl.Println("running start sequence...")

	sl.Println("starting stream, and waiting...")
	if realSession {
		mqtt.Publish(topics_backend.TOPIC_STREAM_START, "")
	}

	wait := time.Duration(streamPreStartMinutes)*time.Minute + time.Second*10
	obs.SetCountdown(
		"Time until start:",
		time.Now().Add(wait),
	)
	time.Sleep(wait)

	// start time
	sl.Println("starting session")
	if realSession {
		mqtt.Publish(topics_backend.TOPIC_SESSION_BEGIN, "PRODUCTION")
	}

	sl.Println("requesting wake")
	mqtt.Publish(topics_firmware.TOPIC_WAKE, "")

	time.Sleep(time.Second * 15)
	// if not awake, abort and error
	sl.Println("checking for awake status")
	sr := events.GetLatestStateReportCopy()
	if sr.Status != machinepb.Status_IDLE_MOVING && sr.Status != machinepb.Status_IDLE_STATIONARY {
		sl.Println("invalid status, erroring")
		mqtt.Publish(topics_backend.TOPIC_SESSION_PAUSE, "")
		err := fmt.Errorf("status was %s after waking but expected idle. Pausing and aborting automation", sr.Status)
		// email for help
		requestSessionIntervention(err)
		return err
	}
	sl.Println("valid status")

	sl.Println("dispensing milk")
	mqtt.Publish(
		topics_firmware.TOPIC_FLUID,
		fmt.Sprintf(
			"%d,%d,%t",
			machinepb.FluidType_FLUID_MILK,
			milkVolume,
			false, // open drain
		),
	)

	time.Sleep(time.Second * 5)
	sl.Println("start sequence done")

	return nil
}

func runEndSequence() error {
	mqtt.Publish(topics_firmware.TOPIC_GOTO_RING_IDLE_POS, "")

	sl.Println("running end sequence")
	sl.Println("draining milk")

	// drain
	mqtt.Publish(
		topics_firmware.TOPIC_FLUID,
		fmt.Sprintf(
			"%d,%d,%t",
			machinepb.FluidType_FLUID_DRAIN,
			milkVolume,
			false, // additional open drain false
		),
	)

	waitForFluidComplete()
	sl.Println("dispensing water")

	mqtt.Publish(
		topics_firmware.TOPIC_FLUID,
		fmt.Sprintf(
			"%d,%d,%t",
			machinepb.FluidType_FLUID_WATER,
			waterVolume,
			false,
		),
	)

	waitForFluidComplete()

	sl.Println("draining water")
	mqtt.Publish(
		topics_firmware.TOPIC_FLUID,
		fmt.Sprintf(
			"%d,%d,%t",
			machinepb.FluidType_FLUID_DRAIN,
			waterVolume,
			false,
		),
	)

	time.Sleep(42 * time.Second)

	sl.Println("shutting down")
	mqtt.Publish(topics_firmware.TOPIC_SHUTDOWN, "")

	waitForSleeping()

	time.Sleep(2 * time.Second)

	sl.Println("ending session")
	mqtt.Publish(topics_backend.TOPIC_SESSION_END, "")

	time.Sleep(2 * time.Minute)

	sl.Println("ending stream")
	mqtt.Publish(topics_backend.TOPIC_STREAM_END, "")

	sl.Println("end sequence done")

	return nil
}

func waitForTOffset(t time.Time, minutesOffset, secondsOffset int) {
	<-time.After(
		time.Until(
			t.
				Add(time.Minute * time.Duration(minutesOffset)).
				Add(time.Second * time.Duration(secondsOffset)),
		),
	)
}

func waitForSleeping() {
	w := events.ConditionWaiter(func(sr *machinepb.StateReport) bool {
		return sr.Status == machinepb.Status_SLEEPING
	})
	<-w
}

func waitForFluidComplete() {
	time.Sleep(2 * time.Second)
	w := events.ConditionWaiter(func(sr *machinepb.StateReport) bool {
		return sr.FluidRequest.Complete
	})
	<-w
}
