package scheduler

import (
	"fmt"
	"time"

	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/asol-protos/go/topics_backend"
	"github.com/gkstretton/asol-protos/go/topics_firmware"
	"github.com/gkstretton/dark/services/goo/actor"
	"github.com/gkstretton/dark/services/goo/events"
	"github.com/gkstretton/dark/services/goo/mqtt"
	"github.com/gkstretton/dark/services/goo/session"
	"github.com/gkstretton/dark/services/goo/twitchapi"
)

type SessionDescriptor struct {
	// how many minutes before session should stream start?
	streamPreStartMinutes  int
	actorDurationMinutes   int
	sessionDurationMinutes int
}

func RunSession(
	d *SessionDescriptor,
	sm *session.SessionManager,
	twitchApi *twitchapi.TwitchApi,
) error {
	mqtt.Publish(topics_backend.TOPIC_STREAM_START, "")
	time.Sleep(time.Duration(d.streamPreStartMinutes) * time.Minute)

	// start time
	beginTime := time.Now()
	mqtt.Publish(topics_backend.TOPIC_SESSION_BEGIN, "PRODUCTION")
	time.Sleep(time.Second * 5)

	mqtt.Publish(topics_firmware.TOPIC_WAKE, "")

	time.Sleep(time.Second * 15)
	// if not awake, abort and error
	sr := events.GetLatestStateReportCopy()
	if sr.Status != machinepb.Status_IDLE_MOVING && sr.Status != machinepb.Status_IDLE_STATIONARY {
		mqtt.Publish(topics_backend.TOPIC_SESSION_PAUSE, "")
		err := fmt.Errorf("status was %s after waking but expected idle. Pausing and aborting automation", sr.Status)
		// email for help
		requestSessionIntervention(err)
		// add timeout to end session after a few hours if no human response
		go sessionTimeout(time.Hour*3, false)
		return err
	}

	mqtt.Publish(
		topics_firmware.TOPIC_FLUID,
		fmt.Sprintf(
			"%d,%d,%t",
			machinepb.FluidType_FLUID_MILK,
			bulkVolume,
			false, // open drain
		),
	)

	time.Sleep(time.Second * 5)

	err := actor.LaunchActor(twitchApi, time.Duration(d.actorDurationMinutes)*time.Minute)
	if err != nil {
		mqtt.Publish(topics_backend.TOPIC_SESSION_PAUSE, "")
		// email for help
		errWrap := fmt.Errorf("actor returned error, unknown situation: %s", err)
		requestSessionIntervention(errWrap)
		// add timeout to end session after a few hours if no human response
		go sessionTimeout(time.Hour*3, true)
		return err
	}

	waitForTOffset := func(t time.Time, minutesOffset, secondsOffset int) {
		<-time.After(
			time.Until(
				t.
					Add(time.Minute * time.Duration(minutesOffset)).
					Add(time.Second * time.Duration(secondsOffset)),
			),
		)
	}

	endTime := beginTime.Add(time.Duration(d.sessionDurationMinutes) * time.Minute)

	waitForTOffset(endTime, -2, -20)

	// drain
	mqtt.Publish(
		topics_firmware.TOPIC_FLUID,
		fmt.Sprintf(
			"%d,%d,%t",
			machinepb.FluidType_FLUID_DRAIN,
			bulkVolume,
			false, // additional open drain false
		),
	)

	waitForTOffset(endTime, -1, -30)

	mqtt.Publish(
		topics_firmware.TOPIC_FLUID,
		fmt.Sprintf(
			"%d,%d,%t",
			machinepb.FluidType_FLUID_WATER,
			bulkVolume,
			true, // open drain
		),
	)

	waitForTOffset(endTime, -0, -30)
	mqtt.Publish(topics_firmware.TOPIC_SHUTDOWN, "")

	waitForTOffset(endTime, 0, 0)
	mqtt.Publish(topics_backend.TOPIC_SESSION_END, "")

	waitForTOffset(endTime, 0, 30)
	mqtt.Publish(topics_backend.TOPIC_STREAM_END, "")
	requestCleaning()
	requestPieceSelection()

	return nil
}
