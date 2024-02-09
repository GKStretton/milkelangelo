package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/gkstretton/dark/services/goo/actor"
	"github.com/gkstretton/dark/services/goo/contentscheduler"
	"github.com/gkstretton/dark/services/goo/email"
	"github.com/gkstretton/dark/services/goo/events"
	"github.com/gkstretton/dark/services/goo/filesystem"
	"github.com/gkstretton/dark/services/goo/keyvalue"
	"github.com/gkstretton/dark/services/goo/livecapture"
	"github.com/gkstretton/dark/services/goo/mqtt"
	"github.com/gkstretton/dark/services/goo/obs"
	"github.com/gkstretton/dark/services/goo/scheduler"
	"github.com/gkstretton/dark/services/goo/server"
	"github.com/gkstretton/dark/services/goo/session"
	"github.com/gkstretton/dark/services/goo/twitchapi"
	"github.com/gkstretton/dark/services/goo/vialprofiles"
)

var (
	test = flag.Bool("test", false, "if true, just run test code")
)

func main() {
	flag.Parse()

	if *test {
		runAdHocTests()
		return
	}

	filesystem.AssertBasePaths()

	mqtt.Start()
	keyvalue.Start()
	email.Start()
	server.Start()

	sm := session.NewSessionManager(false)
	twitchApi := twitchapi.Start()

	actor.Setup()
	events.Start(sm)
	livecapture.Start(sm)
	obs.Start(sm)
	vialprofiles.Start(sm)
	scheduler.Start(sm, twitchApi)
	contentscheduler.Start(sm)

	// Block to prevent early quit
	fmt.Println("finished init, main loop sleeping.")
	for {
		time.Sleep(time.Millisecond * time.Duration(100))
	}
}
