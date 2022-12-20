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

type sessionManager struct {
	s storage
	// pub is the channel that this package can publish to
	pub chan *sessionEvent
	// subs is all the channels listened to by subscribers
	subs []chan *sessionEvent
}

func NewSessionManager(useMemoryStorage bool) *sessionManager {
	sm := &sessionManager{
		s:    newStorage(useMemoryStorage),
		pub:  make(chan *sessionEvent),
		subs: []chan *sessionEvent{},
	}
	go sm.eventDistributor()
	return sm
}

// BeginSession will attempt to begin a new session
func (sm *sessionManager) BeginSession() (*Session, error) {
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
	sm.pub <- &sessionEvent{
		SessionID: session.Id,
		Type:      SESSION_STARTED,
	}
	return session, nil
}

// EndSession will end a session if one is in progress
func (sm *sessionManager) EndSession() (*Session, error) {
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
	sm.pub <- &sessionEvent{
		SessionID: session.Id,
		Type:      SESSION_ENDED,
	}
	return session, nil
}

// GetCurrentSession returns nil, nil if there is no current session
func (sm *sessionManager) GetCurrentSession() (*Session, error) {
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
