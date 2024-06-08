package actor

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/gkstretton/asol-protos/go/topics_backend"
	"github.com/gkstretton/dark/services/goo/mqtt"
	"github.com/gkstretton/dark/services/goo/session"
)

func subscribeToBrokerTopics(sm *session.SessionManager) {
	l.Println("subscribed to topics")
	mqtt.Subscribe(topics_backend.TOPIC_ACTOR_START, func(topic string, payload []byte) {
		l.Println("mqtt start request")
		go func() {
			args := strings.Split(string(payload), ",")
			minutes, err := strconv.Atoi(args[0])
			if err != nil {
				defaultMinutes := 10
				l.Printf("invalid minutes for actor timeout, defaulting to %d: %v\n", defaultMinutes, err)
				minutes = defaultMinutes
			}

			var actorSeed int64
			if len(args) == 2 {
				seed, err := strconv.Atoi(args[1])
				if err != nil {
					l.Printf("invalid seed for actor, defaulting to random: %v\n", err)
					actorSeed = rand.Int63()
				} else {
					actorSeed = int64(seed)
				}
			} else {
				l.Println("no seed provided, defaulting to random")
				actorSeed = rand.Int63()
			}

			err = sm.SetCurrentSessionSeed(actorSeed)
			if err != nil {
				l.Printf("failed to set seed: %v\n", err)
			}

			err = LaunchActor(nil, time.Duration(minutes)*time.Minute, actorSeed, false)
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
