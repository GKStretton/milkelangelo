package session

import (
	"fmt"

	"github.com/gkstretton/dark/services/goo/config"
	"github.com/gkstretton/dark/services/goo/mqtt"
)

type EventType int

const (
	SESSION_STARTED EventType = iota
	SESSION_ENDED
	SESSION_PAUSED
	SESSION_RESUMED
)

type SessionEvent struct {
	SessionID ID
	Type      EventType
}

func (sm *SessionManager) SubscribeToEvents() <-chan *SessionEvent {
	ch := make(chan *SessionEvent)
	sm.subs = append(sm.subs, ch)
	return ch
}

// eventDistributor handles event fan out
func (sm *SessionManager) eventDistributor() {
	fmt.Println("Running eventDistributor")
	for {
		e := <-sm.pub

		// publish event to broker, and to internal channel
		sm.publishToBroker(e)
		for _, sub := range sm.subs {
			// non-blocking send to subscriber
			select {
			case sub <- e:
			default:
			}
		}
	}
}

func (sm *SessionManager) publishToBroker(e *SessionEvent) {
	var topic string
	switch e.Type {
	case SESSION_STARTED:
		topic = config.TOPIC_SESSION_BEGAN
	case SESSION_ENDED:
		topic = config.TOPIC_SESSION_ENDED
	case SESSION_PAUSED:
		topic = config.TOPIC_SESSION_PAUSED
	case SESSION_RESUMED:
		topic = config.TOPIC_SESSION_RESUMED
	default:
		fmt.Printf("unknown event type in publishToBroker: %v\n", e)
		return
	}

	err := mqtt.Publish(topic, fmt.Sprintf("%d", e.SessionID))
	if err != nil {
		fmt.Printf("error publishing session event: %v\n", err)
	}
}

func (sm *SessionManager) subscribeToBrokerTopics() {
	mqtt.Subscribe(config.TOPIC_SESSION_BEGIN, func(topic string, payload []byte) {
		_, err := sm.BeginSession()
		if err != nil {
			fmt.Printf("cannot begin session: %v\n", err)
		}
	})

	mqtt.Subscribe(config.TOPIC_SESSION_END, func(topic string, payload []byte) {
		_, err := sm.EndSession()
		if err != nil {
			fmt.Printf("cannot end session: %v\n", err)
		}
	})

	mqtt.Subscribe(config.TOPIC_SESSION_PAUSE, func(topic string, payload []byte) {
		_, err := sm.PauseSession()
		if err != nil {
			fmt.Printf("cannot pause session: %v\n", err)
		}
	})

	mqtt.Subscribe(config.TOPIC_SESSION_RESUME, func(topic string, payload []byte) {
		_, err := sm.ResumeSession()
		if err != nil {
			fmt.Printf("cannot resume session: %v\n", err)
		}
	})
}
