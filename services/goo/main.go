package main

import (
	"flag"
	"time"

	"github.com/gkstretton/dark/services/goo/events"
	"github.com/gkstretton/dark/services/goo/filesystem"
	"github.com/gkstretton/dark/services/goo/keyvalue"
	"github.com/gkstretton/dark/services/goo/livecapture"
	"github.com/gkstretton/dark/services/goo/mqtt"
	"github.com/gkstretton/dark/services/goo/obs"
	"github.com/gkstretton/dark/services/goo/session"
)

var (
	test = flag.Bool("test", false, "if true, just run test code")
)

func main() {
	flag.Parse()

	if *test {
		Test()
	}

	filesystem.AssertBasePaths()

	mqtt.Start()
	keyvalue.Start()

	sm := session.NewSessionManager(false)

	events.Run(sm)
	livecapture.Run(sm)
	obs.Run(sm)

	// Block to prevent early quit
	for {
		time.Sleep(time.Millisecond * time.Duration(100))
	}
}

func Test() {

}
