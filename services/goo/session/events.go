package session

type EventType int

const (
	SESSION_STARTED EventType = iota
	SESSION_ENDED
)

type sessionEvent struct {
	SessionID ID
	Type      EventType
}

func (sm *sessionManager) GetEventsChan() <-chan *sessionEvent {
	return sm.e
}
