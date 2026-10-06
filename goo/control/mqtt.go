package control

import (
	"fmt"

	"github.com/gkstretton/asol-protos/go/topics_backend"
	"github.com/gkstretton/dark/services/goo/mqtt"
)

// The actor topics are reused for enabling / disabling control.
func subscribeToBrokerTopics() {
	mqtt.Subscribe(topics_backend.TOPIC_ACTOR_START, func(topic string, payload []byte) {
		l.Println("mqtt enable request")
		SetEnabled(true)
	})
	mqtt.Subscribe(topics_backend.TOPIC_ACTOR_STOP, func(topic string, payload []byte) {
		l.Println("mqtt disable request")
		SetEnabled(false)
	})
	mqtt.Subscribe(topics_backend.TOPIC_ACTOR_STATUS_GET, func(topic string, payload []byte) {
		publishStatus()
	})

	publishStatus()
}

func publishStatus() {
	mqtt.Publish(topics_backend.TOPIC_ACTOR_STATUS_RESP, fmt.Sprintf("%t", IsEnabled()))
}
