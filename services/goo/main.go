package main

import (
	"flag"
	"time"

	"github.com/gkstretton/dark/services/goo/email"
	"github.com/gkstretton/dark/services/goo/events"
	"github.com/gkstretton/dark/services/goo/filesystem"
	"github.com/gkstretton/dark/services/goo/keyvalue"
	"github.com/gkstretton/dark/services/goo/livecapture"
	"github.com/gkstretton/dark/services/goo/mqtt"
	"github.com/gkstretton/dark/services/goo/obs"
	"github.com/gkstretton/dark/services/goo/scheduler"
	"github.com/gkstretton/dark/services/goo/session"
	"github.com/gkstretton/dark/services/goo/vialprofiles"
)

var (
	test = flag.Bool("test", false, "if true, just run test code")
)

func main() {
	flag.Parse()

	if *test {
		scheduler.TestEntry()
		return
	}

	filesystem.AssertBasePaths()

	mqtt.Start()
	keyvalue.Start()
	email.Start()

	sm := session.NewSessionManager(false)

	events.Start(sm)
	livecapture.Start(sm)
	obs.Start(sm)
	vialprofiles.Start(sm)
	scheduler.Start()

	// Block to prevent early quit
	for {
		time.Sleep(time.Millisecond * time.Duration(100))
	}
}
