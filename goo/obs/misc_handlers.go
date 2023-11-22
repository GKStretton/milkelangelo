package obs

import (
	"fmt"

	"github.com/gkstretton/asol-protos/go/topics_backend"
	"github.com/gkstretton/dark/services/goo/mqtt"
)

func onStreamStart() {
	fmt.Printf("Closing blind, monitor on...\n")
	err := mqtt.Publish(topics_backend.TOPIC_MONITOR_ON, "")
	if err != nil {
		fmt.Printf("error publishing monitor on: %v\n", err)
	}
	err = mqtt.Publish(topics_backend.TOPIC_CLOSE_BLIND, "")
	if err != nil {
		fmt.Printf("error publishing close blind: %v\n", err)
	}
}

func onStreamEnd() {
	fmt.Printf("Opening blind, monitor off...\n")
	err := mqtt.Publish(topics_backend.TOPIC_MONITOR_OFF, "")
	if err != nil {
		fmt.Printf("error publishing monitor on: %v\n", err)
	}
	err = mqtt.Publish(topics_backend.TOPIC_OPEN_BLIND, "")
	if err != nil {
		fmt.Printf("error publishing open blind: %v\n", err)
	}
}
