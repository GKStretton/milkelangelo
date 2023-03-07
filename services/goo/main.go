package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/dark/services/goo/config"
	"github.com/gkstretton/dark/services/goo/events"
	"github.com/gkstretton/dark/services/goo/filesystem"
	"github.com/gkstretton/dark/services/goo/keyvalue"
	"github.com/gkstretton/dark/services/goo/livecapture"
	"github.com/gkstretton/dark/services/goo/mqtt"
	"github.com/gkstretton/dark/services/goo/obs"
	"github.com/gkstretton/dark/services/goo/session"
	"google.golang.org/protobuf/encoding/protojson"
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
	mqtt.Start()

	sr := &machinepb.StateReport{
		Mode:              machinepb.Mode_AUTONOMOUS,
		Status:            machinepb.Status_SLEEPING,
		PipetteState:      &machinepb.PipetteState{},
		CollectionRequest: &machinepb.CollectionRequest{},
		MovementDetails:   &machinepb.MovementDetails{},
	}
	m := protojson.MarshalOptions{
		Multiline:       true,
		UseProtoNames:   true,
		Indent:          "\t",
		EmitUnpopulated: true,
	}
	b, err := m.Marshal(sr)
	fmt.Printf("%s\n%v\n", string(b), err)

	err = mqtt.Publish(config.TOPIC_STATE_REPORT_JSON, string(b))
	if err != nil {
		fmt.Printf("error publishing json state report: %v\n", err)
		return
	}
}
