package events

import (
	"fmt"

	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/dark/services/goo/mqtt"
	"google.golang.org/protobuf/proto"
)

func Run() {
	mqtt.Subscribe("mega/state-report", func(topic string, payload []byte) {
		fmt.Printf("Received machine state report: '%v', ' %s '\n", payload, string(payload))

		sr := &machinepb.StateReport{}
		err := proto.Unmarshal(payload, sr)
		if err != nil {
			fmt.Printf("error unmarshalling state report: %v\n", err)
			return
		}

		fmt.Printf("%+v\n", sr)
	})

	//todo: add a way to fake the state reports, and continue building
}
