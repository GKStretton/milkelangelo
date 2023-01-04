package main

import (
	"flag"
	"time"

	"github.com/gkstretton/dark/services/goo/events"
	"github.com/gkstretton/dark/services/goo/filesystem"
	"github.com/gkstretton/dark/services/goo/keyvalue"
	"github.com/gkstretton/dark/services/goo/livecapture"
	"github.com/gkstretton/dark/services/goo/mqtt"
	"github.com/gkstretton/dark/services/goo/session"
)

func main() {
	flag.Parse()

	filesystem.AssertBasePaths()

	mqtt.Start()
	keyvalue.Start()

	sm := session.NewSessionManager(false)

	events.Run(sm)
	livecapture.Run(sm)

	// Block to prevent early quit
	for {
		time.Sleep(time.Millisecond * time.Duration(100))
	}
}
