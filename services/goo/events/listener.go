package events

import (
	"fmt"

	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/dark/services/goo/mqtt"
	"google.golang.org/protobuf/proto"
)

func Run() {
	mqtt.Subscribe("mega/state-report", func(topic string, payload []byte) {
		fmt.Println("Received machine state report")

		pr := &machinepb.PingResponse{}
		err := proto.Unmarshal(payload, pr)
		if err != nil {
			fmt.Printf("error unmarshalling state report: %v\n", err)
		}

		fmt.Printf("\tNumber: %d\n", pr.Number)
	})
}
