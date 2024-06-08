package session

import (
	"fmt"
	"sync"

	"github.com/gkstretton/asol-protos/go/topics_backend"
	"github.com/gkstretton/dark/services/goo/mqtt"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var eventLock sync.Mutex

type EventType int

const (
	SESSION_STARTED EventType = iota
	SESSION_ENDED
	SESSION_PAUSED
	SESSION_RESUMED
)

const EVENT_CHAN_BUFFER = 10

type SessionEvent struct {
	SessionID ID
	Type      EventType
}

func (sm *SessionManager) SubscribeToEvents() <-chan *SessionEvent {
	ch := make(chan *SessionEvent, EVENT_CHAN_BUFFER)

	eventLock.Lock()
	defer eventLock.Unlock()

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
		eventLock.Lock()
		for _, sub := range sm.subs {
			// non-blocking send to subscriber
			select {
			case sub <- e:
			default:
			}
		}
		eventLock.Unlock()
	}
}

func (sm *SessionManager) publishToBroker(e *SessionEvent) {
	var topic string
	switch e.Type {
	case SESSION_STARTED:
		topic = topics_backend.TOPIC_SESSION_BEGAN
	case SESSION_ENDED:
		topic = topics_backend.TOPIC_SESSION_ENDED
	case SESSION_PAUSED:
		topic = topics_backend.TOPIC_SESSION_PAUSED
	case SESSION_RESUMED:
		topic = topics_backend.TOPIC_SESSION_RESUMED
	default:
		fmt.Printf("unknown event type in publishToBroker: %v\n", e)
		return
	}

	err := mqtt.Publish(topic, fmt.Sprintf("%d", e.SessionID))
	if err != nil {
		fmt.Printf("error publishing session event: %v\n", err)
	}

	sm.publishSessionStatus()
}

func (sm *SessionManager) publishSessionStatus() {
	session, _ := sm.GetLatestSession()
	s := session.ToProto()

	// protobuf
	b, err := proto.Marshal(s)
	if err != nil {
		fmt.Printf("error marshalling session status as protobuf: %v\n", err)
	}
	if err = mqtt.Publish(topics_backend.TOPIC_SESSION_STATUS_RESP_RAW, b); err != nil {
		fmt.Printf("error publishing session status: %v\n", err)
	}

	// json
	m := protojson.MarshalOptions{
		Multiline:       true,
		UseProtoNames:   true,
		Indent:          "\t",
		EmitUnpopulated: true,
	}
	j, err := m.Marshal(s)
	if err != nil {
		fmt.Printf("error marshalling session status to json: %v\n", err)
	}
	if err = mqtt.Publish(topics_backend.TOPIC_SESSION_STATUS_RESP_JSON, j); err != nil {
		fmt.Printf("error publishing session status: %v\n", err)
	}
}

func (sm *SessionManager) subscribeToBrokerTopics() {
	mqtt.Subscribe(topics_backend.TOPIC_SESSION_BEGIN, func(topic string, payload []byte) {
		var production bool
		if string(payload) == "PRODUCTION" {
			production = true
		}
		fmt.Printf("received mqtt request to begin session, production=%t\n", production)
		go func() {
			_, err := sm.BeginSession(production)
			if err != nil {
				fmt.Printf("cannot begin session: %v\n", err)
			}
		}()
	})

	mqtt.Subscribe(topics_backend.TOPIC_SESSION_END, func(topic string, payload []byte) {
		fmt.Println("received mqtt request to end session")
		go func() {
			_, err := sm.EndSession()
			if err != nil {
				fmt.Printf("cannot end session: %v\n", err)
			}
		}()
	})

	mqtt.Subscribe(topics_backend.TOPIC_SESSION_PAUSE, func(topic string, payload []byte) {
		fmt.Println("received mqtt request to pause session")
		go func() {
			_, err := sm.PauseSession()
			if err != nil {
				fmt.Printf("cannot pause session: %v\n", err)
			}
		}()
	})

	mqtt.Subscribe(topics_backend.TOPIC_SESSION_RESUME, func(topic string, payload []byte) {
		fmt.Println("received mqtt request to resume session")
		go func() {
			_, err := sm.ResumeSession()
			if err != nil {
				fmt.Printf("cannot resume session: %v\n", err)
			}
		}()
	})

	mqtt.Subscribe(topics_backend.TOPIC_SESSION_STATUS_GET, func(topic string, payload []byte) {
		go sm.publishSessionStatus()
	})
}
