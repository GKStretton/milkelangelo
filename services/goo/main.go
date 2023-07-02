package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/gkstretton/dark/services/goo/email"
	"github.com/gkstretton/dark/services/goo/events"
	"github.com/gkstretton/dark/services/goo/filesystem"
	"github.com/gkstretton/dark/services/goo/keyvalue"
	"github.com/gkstretton/dark/services/goo/livecapture"
	"github.com/gkstretton/dark/services/goo/mqtt"
	"github.com/gkstretton/dark/services/goo/obs"
	"github.com/gkstretton/dark/services/goo/session"
	"github.com/gkstretton/dark/services/goo/vialprofiles"
)

var (
	test = flag.Bool("test", false, "if true, just run test code")
)

func main() {
	flag.Parse()

	if *test {
		Test()
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

	// Block to prevent early quit
	for {
		time.Sleep(time.Millisecond * time.Duration(100))
	}
}

func Test() {
	micros := uint64(1677327218577344)
	t := time.UnixMicro(int64(micros))
	str := t.Format("2006-03-02 15:04:05.000000")
	fmt.Println(micros)
	fmt.Println(str)
}
