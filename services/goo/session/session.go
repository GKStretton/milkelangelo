package session

import (
	"fmt"

	"github.com/gkstretton/dark/services/goo/util"
)

type ID uint64

type Session struct {
	Id         ID
	Complete   bool
	Production bool
}

type SessionMatcher struct {
	Id         *ID
	Complete   *bool
	Production *bool
}

type SessionManager struct {
	s storage
	// pub is the channel that this package can publish to
	pub chan *SessionEvent
	// subs is all the channels listened to by subscribers
	subs []chan *SessionEvent
}

func NewSessionManager(useMemoryStorage bool) *SessionManager {
	sm := &SessionManager{
		s:    newStorage(useMemoryStorage),
		pub:  make(chan *SessionEvent),
		subs: []chan *SessionEvent{},
	}
	go sm.eventDistributor()

	sm.subscribeToBrokerTopics()

	return sm
}

// BeginSession will attempt to begin a new session
func (sm *SessionManager) BeginSession() (*Session, error) {
	current, err := sm.GetCurrentSession()
	if err != nil {
		return nil, fmt.Errorf("failed to check if session is in progress: %v", err)
	}
	if current != nil {
		return nil, fmt.Errorf("session already in progress")
	}

	// okay to start a session
	session := &Session{
		Complete:   false,
		Production: false,
	}
	session, err = sm.s.createSession(session)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %v", err)
	}
	sm.pub <- &SessionEvent{
		SessionID: session.Id,
		Type:      SESSION_STARTED,
	}
	fmt.Printf("Began session %d\n", session.Id)
	return session, nil
}

// EndSession will end a session if one is in progress
func (sm *SessionManager) EndSession() (*Session, error) {
	session, err := sm.GetCurrentSession()
	if err != nil {
		return nil, fmt.Errorf("failed getting current session: %v", err)
	}
	if session == nil {
		return nil, fmt.Errorf("no session in progress")
	}
	session.Complete = true
	session, err = sm.s.updateSession(session)
	if err != nil {
		return nil, fmt.Errorf("failed to update session: %v", err)
	}
	sm.pub <- &SessionEvent{
		SessionID: session.Id,
		Type:      SESSION_ENDED,
	}
	fmt.Printf("Ended session %d\n", session.Id)
	return session, nil
}

// GetCurrentSession returns nil, nil if there is no current session
func (sm *SessionManager) GetCurrentSession() (*Session, error) {
	incompleteSessions, err := sm.s.matchSession(&SessionMatcher{
		Complete: util.Ptr(false),
	})
	if err != nil {
		return nil, err
	}
	if len(incompleteSessions) > 1 {
		return nil, fmt.Errorf("%d sessions in progress, expected <= 1", len(incompleteSessions))
	}

	if len(incompleteSessions) == 1 {
		return incompleteSessions[0], nil
	}

	return nil, nil
}
