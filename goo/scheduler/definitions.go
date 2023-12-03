package scheduler

import (
	"fmt"
	"time"

	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/asol-protos/go/topics_backend"
	"github.com/gkstretton/asol-protos/go/topics_firmware"
	"github.com/gkstretton/dark/services/goo/actor"
	"github.com/gkstretton/dark/services/goo/mqtt"
	"github.com/gkstretton/dark/services/goo/session"
	"github.com/gkstretton/dark/services/goo/twitchapi"
)

const bulkVolume = 200 //ml of milk

// time the session will begin UTC
var mainSessionStartTime = RecurringTime{
	day:    time.Saturday,
	hour:   18,
	minute: 30,
	second: 0,
}

// time the session will end UTC
var mainSessionEndTime = RecurringTime{
	day:    time.Saturday,
	hour:   19,
	minute: 30,
	second: 0,
}

// defineSchedule works by launching go routines watching for the specified
// time, to trigger the stated action.
func defineSchedule(sm *session.SessionManager, twitchApi *twitchapi.TwitchApi) {
	go scheduleWatcher(&Schedule{
		name:    "FRIDGE_ON",
		enabled: true,
		function: func() {
			mqtt.Publish(topics_backend.TOPIC_FRIDGE_SWITCH, topics_backend.PAYLOAD_SMART_SWITCH_ON)
		},
		recurringTime: mainSessionStartTime,
		hourOffset:    -7,
		minuteOffset:  -30,
	})

	// ***********
	// SESSION START
	// ***********

	go scheduleWatcher(&Schedule{
		name:    "START_STREAM",
		enabled: true,
		function: func() {
			mqtt.Publish(topics_backend.TOPIC_STREAM_START, "")
		},
		recurringTime: mainSessionStartTime,
		minuteOffset:  -5,
	})

	go scheduleWatcher(&Schedule{
		name:    "BEGIN_PROD_SESSION",
		enabled: true,
		function: func() {
			mqtt.Publish(topics_backend.TOPIC_SESSION_BEGIN, "PRODUCTION")
		},
		recurringTime: mainSessionStartTime,
	})

	go scheduleWatcher(&Schedule{
		name:    "WAKE_ROBOT",
		enabled: true,
		function: func() {
			mqtt.Publish(topics_firmware.TOPIC_WAKE, "")
		},
		recurringTime: mainSessionStartTime,
		minuteOffset:  0,
		secondOffset:  10,
	})

	go scheduleWatcher(&Schedule{
		name:    "DISPENSE_MILK",
		enabled: true,
		function: func() {
			mqtt.Publish(
				topics_firmware.TOPIC_FLUID,
				fmt.Sprintf(
					"%d,%d,%t",
					machinepb.FluidType_FLUID_MILK,
					bulkVolume,
					false, // open drain
				),
			)
		},
		recurringTime: mainSessionStartTime,
		minuteOffset:  0,
		secondOffset:  20,
	})

	// ***********
	// AUTOMATED ACTOR
	// ***********

	go scheduleWatcher(&Schedule{
		name:    "LAUNCH_ACTOR",
		enabled: true,
		function: func() {
			actor.LaunchActor(twitchApi, 5*time.Second)
		},
		recurringTime: mainSessionStartTime,
		minuteOffset:  0,
		secondOffset:  25,
	})

	// ***********
	// SESSION END
	// ***********

	go scheduleWatcher(&Schedule{
		name:    "DRAIN",
		enabled: true,
		function: func() {
			mqtt.Publish(
				topics_firmware.TOPIC_FLUID,
				fmt.Sprintf(
					"%d,%d,%t",
					machinepb.FluidType_FLUID_DRAIN,
					bulkVolume,
					false, // open drain
				),
			)
		},
		recurringTime: mainSessionEndTime,
		minuteOffset:  -2,
		secondOffset:  -20,
	})

	go scheduleWatcher(&Schedule{
		name:    "RINSE",
		enabled: true,
		function: func() {
			mqtt.Publish(
				topics_firmware.TOPIC_FLUID,
				fmt.Sprintf(
					"%d,%d,%t",
					machinepb.FluidType_FLUID_WATER,
					bulkVolume,
					true, // open drain
				),
			)
		},
		recurringTime: mainSessionEndTime,
		minuteOffset:  -1,
		secondOffset:  -30,
	})

	go scheduleWatcher(&Schedule{
		name:    "SHUTDOWN_ROBOT",
		enabled: true,
		function: func() {
			mqtt.Publish(topics_firmware.TOPIC_SHUTDOWN, "")
		},
		recurringTime: mainSessionEndTime,
		secondOffset:  -30,
	})

	go scheduleWatcher(&Schedule{
		name:    "END_SESSION",
		enabled: true,
		function: func() {
			mqtt.Publish(topics_backend.TOPIC_SESSION_END, "")
		},
		recurringTime: mainSessionEndTime,
	})

	go scheduleWatcher(&Schedule{
		name:    "END_STREAM",
		enabled: true,
		function: func() {
			mqtt.Publish(topics_backend.TOPIC_STREAM_END, "")
		},
		recurringTime: mainSessionEndTime,
		secondOffset:  30,
	})

	go scheduleWatcher(&Schedule{
		name:    "FRIDGE_OFF",
		enabled: true,
		function: func() {
			mqtt.Publish(topics_backend.TOPIC_FRIDGE_SWITCH, topics_backend.PAYLOAD_SMART_SWITCH_OFF)
		},
		recurringTime: mainSessionEndTime,
		minuteOffset:  1,
	})
}
