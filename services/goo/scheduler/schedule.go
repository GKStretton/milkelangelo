package scheduler

import (
	"fmt"
	"time"
)

type Schedule struct {
	name     string
	function func()
	day      time.Weekday
	hour     int
	minute   int
	second   int
}

func (s *Schedule) nextRunTime() time.Time {
	// schedule items are defined in utc
	now := time.Now().UTC()
	// work out the time
	nextRunTime := time.Date(now.Year(), now.Month(), now.Day(), s.hour, s.minute, s.second, 0, now.Location())
	// fix the day
	nextRun := nextRunTime.AddDate(0, 0, int(s.day)-int(now.Weekday()))
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
		// execute
		fmt.Printf("[schedule] executing %s\n", s.name)
		s.function()
	}
}

func defineSchedules() {
	go scheduleWatcher(&Schedule{
		name: "START_STREAM",
		function: func() {
			fmt.Println("test START_STREAM")
		},
		day:    time.Sunday,
		hour:   20,
		minute: 12,
		second: 10,
	})
}
