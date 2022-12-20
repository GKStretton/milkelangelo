package session

import "fmt"

type EventType int

const (
	SESSION_STARTED EventType = iota
	SESSION_ENDED
)

type sessionEvent struct {
	SessionID ID
	Type      EventType
}

func (sm *sessionManager) SubscribeToEvents() (<-chan *sessionEvent, error) {
	ch := make(chan *sessionEvent)
	sm.subs = append(sm.subs, ch)
	return ch, nil
}

// eventDistributor handles event fan out
func (sm *sessionManager) eventDistributor() {
	fmt.Println("Running eventDistributor")
	for {
		e := <-sm.pub
		for _, sub := range sm.subs {
			// non-blocking send to subscriber
			select {
			case sub <- e:
			default:
			}
		}
	}
}
