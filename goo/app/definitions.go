package app

import (
	"fmt"
	"time"

	"github.com/gkstretton/asol-protos/go/topics_backend"
	"github.com/gkstretton/dark/services/goo/ebsinterface"
	"github.com/gkstretton/dark/services/goo/keyvalue"
	"github.com/gkstretton/dark/services/goo/mqtt"
	"github.com/gkstretton/dark/services/goo/session"
	"github.com/gkstretton/dark/services/goo/twitchapi"
)

const (
	//ml of milk
	milkVolume                    = 200
	waterVolume                   = 300
	streamPreStartMinutes         = 10
	defaultSessionDurationMinutes = 50
)

// time the session will begin UTC
var mainSessionStartTime = RecurringTime{
	day:    time.Sunday,
	hour:   17,
	minute: 0,
	second: 0,
}

// defineSchedule works by launching go routines watching for the specified
// time, to trigger the stated action.
func defineSchedule(sm *session.SessionManager, twitchApi *twitchapi.TwitchApi, ebsApi ebsinterface.EbsApi) {
	go scheduleWatcher(&Schedule{
		name:    "REMINDER",
		enabled: true,
		function: func() {
			s := readOneTimeSettings()

			// remind routine operator
			sendReminder(s.skip)
		},
		recurringTime: mainSessionStartTime,
		hourOffset:    -11,
		minuteOffset:  0,
	})

	go scheduleWatcher(&Schedule{
		name:    "REQUEST_MILK",
		enabled: true,
		function: func() {
			s := readOneTimeSettings()

			if s.skip {
				fmt.Println("one time settings skip flag set, skipping milk request!")
				return
			}

			// disabling fridge to reduce convection
			// mqtt.Publish(topics_backend.TOPIC_FRIDGE_SWITCH, topics_backend.PAYLOAD_SMART_SWITCH_ON)

			// notify routine operator to fill with milk
			requestFridgeMilk()
		},
		recurringTime: mainSessionStartTime,
		hourOffset:    -3,
		minuteOffset:  0,
	})

	go scheduleWatcher(&Schedule{
		name:    "RUN_SESSION",
		enabled: true,
		function: func() {
			s := readOneTimeSettings()
			resetOneTimeSettings()

			if s.skip {
				fmt.Println("one time settings skip flag set, skipping session!")
				return
			}

			err := RunSession(
				&SessionDescriptor{
					streamPreStartMinutes:  streamPreStartMinutes,
					actorDurationMinutes:   actorDurationMins,
					sessionDurationMinutes: defaultSessionDurationMinutes,
					runActor:               !s.disableActor,
				},
				sm, twitchApi, ebsApi,
			)
			if err != nil {
				fmt.Println(err)
			}
		},
		recurringTime: mainSessionStartTime,
		minuteOffset:  -streamPreStartMinutes,
	})

	go scheduleWatcher(&Schedule{
		name:    "FRIDGE_OFF",
		enabled: false,
		function: func() {
			mqtt.Publish(topics_backend.TOPIC_FRIDGE_SWITCH, topics_backend.PAYLOAD_SMART_SWITCH_OFF)
		},
		recurringTime: mainSessionStartTime,
		hourOffset:    1, // 1 hours after session start
		minuteOffset:  0,
	})
}

type oneTimeSettings struct {
	skip         bool
	disableActor bool
}

// readOneTimeSettings checks for one time settings and returns them, resetting
// the flags to false afterwards.
func readOneTimeSettings() *oneTimeSettings {
	skip := keyvalue.GetBool(topics_backend.KV_SCHEDULED_SESSION_FLAG_SKIP)
	disableActor := keyvalue.GetBool(topics_backend.KV_SCHEDULED_SESSION_FLAG_DISABLE_ACTOR)

	return &oneTimeSettings{
		skip:         skip,
		disableActor: disableActor,
	}
}

func resetOneTimeSettings() {
	keyvalue.SetBool(topics_backend.KV_SCHEDULED_SESSION_FLAG_SKIP, false)
	keyvalue.SetBool(topics_backend.KV_SCHEDULED_SESSION_FLAG_DISABLE_ACTOR, false)
}
