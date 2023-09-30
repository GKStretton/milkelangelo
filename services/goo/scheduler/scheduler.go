package scheduler

import (
	"time"
)

func Start() {
	defineSchedules()
}

func TestEntry() {
	Start()
	for {
		time.Sleep(time.Second)
	}
}
