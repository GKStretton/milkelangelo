package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/dark/services/goo/events"
	"github.com/gkstretton/dark/services/goo/filesystem"
	"github.com/gkstretton/dark/services/goo/keyvalue"
	"github.com/gkstretton/dark/services/goo/livecapture"
	"github.com/gkstretton/dark/services/goo/mqtt"
	"github.com/gkstretton/dark/services/goo/session"
	"github.com/gkstretton/dark/services/goo/util/protoyaml"
	"google.golang.org/protobuf/encoding/protojson"
)

var (
	test = flag.Bool("test", false, "if true, just run test code")
)

func main() {
	flag.Parse()

	if *test {
		testFunc()
		return
	}

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

func testFunc() {
	sr := &machinepb.StateReportList{
		StateReports: []*machinepb.StateReport{
			{
				TimestampUnixMicros: 100,
				Mode:                machinepb.Mode_AUTONOMOUS,
				Status:              machinepb.Status_CLEANING_BOWL,
			},
		},
	}
	out, err := protojson.Marshal(sr)
	fmt.Println(err)
	fmt.Println(string(out))
	out, err = protoyaml.Marshal(sr)
	fmt.Println(err)
	fmt.Println(string(out))
}
