package scheduler

import (
	"fmt"
	"time"

	"github.com/gkstretton/dark/services/goo/keyvalue"
	"github.com/gkstretton/dark/services/goo/session"
	"github.com/gkstretton/dark/services/goo/twitchapi"
)

type RecurringTime struct {
	day    time.Weekday
	hour   int
	minute int
	second int
}

// fmtLocal gets just the hh:mm in local time for the maintainer (London)
func (r RecurringTime) fmtLocal() string {
	now := time.Now().UTC()
	d := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		mainSessionStartTime.hour,
		mainSessionStartTime.minute,
		mainSessionStartTime.second,
		0,
		now.Location(),
	)
	localTz, err := time.LoadLocation("Europe/London")
	if err != nil {
		fmt.Printf("error loading tz location: %s\n", err)
		return fmt.Sprintf("%d:%d", d.Hour(), d.Minute())
	}
	d = d.In(localTz)
	return fmt.Sprintf("%d:%d", d.Hour(), d.Minute())
}

type Schedule struct {
	name          string
	enabled       bool
	function      func()
	recurringTime RecurringTime
	// this long before or after the recurring time
	hourOffset   int
	minuteOffset int
	secondOffset int
}

func (s *Schedule) nextRunTime() time.Time {
	// schedule items are defined in utc
	now := time.Now().UTC()
	// work out the time
	recurringTime := s.recurringTime
	nextRunTime := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		recurringTime.hour,
		recurringTime.minute,
		recurringTime.second,
		0,
		now.Location(),
	)
	// offset the time
	nextRunTime = nextRunTime.Add(time.Hour * time.Duration(s.hourOffset))
	nextRunTime = nextRunTime.Add(time.Minute * time.Duration(s.minuteOffset))
	nextRunTime = nextRunTime.Add(time.Second * time.Duration(s.secondOffset))

	// fix the day
	nextRun := nextRunTime.AddDate(0, 0, int(recurringTime.day)-int(now.Weekday()))
	if nextRun.Before(now) {
		nextRun = nextRun.AddDate(0, 0, 7)
	}
	return nextRun
}

func scheduleWatcher(s *Schedule) {
	for {
		// wait until next s.
		nextRun := s.nextRunTime()
		timeUntil := time.Until(nextRun)
		fmt.Printf("[schedule] scheduling %s for %s (%s)\n", s.name, nextRun, timeUntil)
		<-time.After(timeUntil)

		if !s.enabled || !keyvalue.GetBool("ENABLE_SCHEDULER") {
			fmt.Printf("[schedule] skipping %s, not enabled\n", s.name)
			continue
		}

		// execute
		fmt.Printf("[schedule] executing %s\n", s.name)
		s.function()
	}
}

func Start(sm *session.SessionManager, twitchApi *twitchapi.TwitchApi) {
	fmt.Printf("Starting scheduler\n")
	registerHandlers(sm, twitchApi)
	defineSchedule(sm, twitchApi)
}
