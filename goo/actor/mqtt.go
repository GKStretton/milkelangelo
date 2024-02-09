package actor

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gkstretton/asol-protos/go/topics_backend"
	"github.com/gkstretton/dark/services/goo/mqtt"
)

func subscribeToBrokerTopics() {
	l.Println("subscribed to topics")
	mqtt.Subscribe(topics_backend.TOPIC_ACTOR_START, func(topic string, payload []byte) {
		l.Println("mqtt start request")
		go func() {
			minutes, err := strconv.Atoi(string(payload))
			if err != nil {
				defaultMinutes := 10
				l.Printf("invalid minutes for actor timeout, defaulting to %d: %v\n", defaultMinutes, err)
				minutes = defaultMinutes
			}

			err = LaunchActor(nil, time.Duration(minutes)*time.Minute)
			if err != nil {
				l.Printf("error when running actor (mqtt trigger): %v\n", err)
			}
			l.Println("mqtt-triggered actor completed")
		}()
	})
	mqtt.Subscribe(topics_backend.TOPIC_ACTOR_STOP, func(topic string, payload []byte) {
		l.Println("mqtt stop request")
		select {
		case exitCh <- struct{}{}:
			l.Println("send actor stop request")
		default:
			l.Println("didn't send actor stop request, channel already contains stop")
		}
	})
	mqtt.Subscribe(topics_backend.TOPIC_ACTOR_STATUS_GET, func(topic string, payload []byte) {
		l.Println("mqtt get request")
		mqtt.Publish(topics_backend.TOPIC_ACTOR_STATUS_RESP, fmt.Sprintf("%t", lock.Get()))
	})

	mqtt.Publish(topics_backend.TOPIC_ACTOR_STATUS_RESP, fmt.Sprintf("%t", lock.Get()))
}
