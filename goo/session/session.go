package session

import (
	"fmt"
	"sync"

	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/asol-protos/go/topics_firmware"
	"github.com/gkstretton/dark/services/goo/filesystem"
	"github.com/gkstretton/dark/services/goo/mqtt"
)

type ID uint64

type Session struct {
	Id           ID    `yaml:"id"`
	Paused       bool  `yaml:"paused"`
	Complete     bool  `yaml:"complete"`
	Production   bool  `yaml:"production"`
	ProductionId ID    `yaml:"production_id"`
	Seed         int64 `yaml:"seed"`
}

func (s *Session) copy() *Session {
	if s == nil {
		return nil
	}
	return &Session{
		Id:           s.Id,
		Paused:       s.Paused,
		Complete:     s.Complete,
		Production:   s.Production,
		ProductionId: s.ProductionId,
		Seed:         s.Seed,
	}
}

func (s *Session) ToProto() *machinepb.SessionStatus {
	if s == nil {
		return nil
	}

	return &machinepb.SessionStatus{
		Id:           uint64(s.Id),
		Paused:       s.Paused,
		Complete:     s.Complete,
		Production:   s.Production,
		ProductionId: uint64(s.ProductionId),
	}
}

type SessionMatcher struct {
	Id           *ID
	Paused       *bool
	Complete     *bool
	Production   *bool
	ProductionId *ID
}

type SessionManager struct {
	s storage
	// pub is the channel that this package can publish to
	pub chan *SessionEvent
	// subs is all the channels listened to by subscribers
	subs []chan *SessionEvent
	// True if sessions are stored in ram
	inMemory           bool
	latestSessionCache *Session
	lock               *sync.Mutex
}

func NewSessionManager(useMemoryStorage bool) *SessionManager {
	sm := &SessionManager{
		s:        newStorage(useMemoryStorage),
		pub:      make(chan *SessionEvent),
		subs:     []chan *SessionEvent{},
		inMemory: useMemoryStorage,
		lock:     &sync.Mutex{},
	}
	go sm.eventDistributor()

	sm.subscribeToBrokerTopics()

	return sm
}

func (sm *SessionManager) SetCurrentSessionSeed(seed int64) error {
	latest, err := sm.GetLatestSession()
	if err != nil {
		return fmt.Errorf("failed to check if session is in progress: %v", err)
	}
	if latest == nil || latest.Complete {
		return fmt.Errorf("session not in progress")
	}
	latest.Seed = seed
	_, err = sm.s.updateSession(latest)
	if err != nil {
		return fmt.Errorf("failed to update session: %v", err)
	}
	sm.clearLatestCache()
	return nil
}

// BeginSession will attempt to begin a new session
func (sm *SessionManager) BeginSession(production bool) (*Session, error) {
	latest, err := sm.GetLatestSession()
	if err != nil {
		return nil, fmt.Errorf("failed to check if session is in progress: %v", err)
	}
	if latest != nil && !latest.Complete {
		if !latest.Paused {
			return nil, fmt.Errorf("session already in progress and unpaused")
		}
		fmt.Println("Paused session in progress, resuming instead of starting")
		return sm.ResumeSession()
	}

	// okay to start a session
	session := &Session{
		Complete:   false,
		Production: production,
	}
	// Save created session somewhere, db or filesystem
	session, err = sm.s.createSession(session)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %v", err)
	}
	sm.clearLatestCache()

	if !sm.inMemory {
		// Create session folder for content etc.
		err = filesystem.InitSessionContent(uint64(session.Id), uint64(session.ProductionId))
		if err != nil {
			sm.s.deleteSession(session.Id)
			return nil, fmt.Errorf("failed to InitSession in filesystem: %v", err)
		}
	}

	// Notify listeners
	sm.pub <- &SessionEvent{
		SessionID: session.Id,
		Type:      SESSION_STARTED,
	}
	requestStateReport()
	fmt.Printf("Began session %d\n", session.Id)

	return session, nil
}

// ResumeSession will resume a paused in-progress session
func (sm *SessionManager) ResumeSession() (*Session, error) {
	latest, err := sm.GetLatestSession()
	if err != nil {
		return nil, fmt.Errorf("failed to check if session is in progress: %v", err)
	}
	if latest == nil || latest.Complete {
		return nil, fmt.Errorf("no session in progress, cannot resume")
	}

	if !latest.Paused {
		return latest, fmt.Errorf("session not paused")
	}
	latest.Paused = false
	latest, err = sm.s.updateSession(latest)
	if err != nil {
		return latest, fmt.Errorf("failed to resume: %v", err)
	}
	sm.clearLatestCache()

	sm.pub <- &SessionEvent{
		SessionID: latest.Id,
		Type:      SESSION_RESUMED,
	}
	requestStateReport()
	fmt.Printf("Resumed session %d\n", latest.Id)

	return latest, nil
}

// PauseSession will pause a current session
func (sm *SessionManager) PauseSession() (*Session, error) {
	latest, err := sm.GetLatestSession()
	if err != nil {
		return nil, fmt.Errorf("failed to check if session is in progress: %v", err)
	}
	if latest == nil || latest.Complete {
		return nil, fmt.Errorf("no session in progress, cannot pause")
	}

	if latest.Paused {
		return latest, fmt.Errorf("session already paused")
	}
	latest.Paused = true
	latest, err = sm.s.updateSession(latest)
	if err != nil {
		return latest, fmt.Errorf("failed to pause: %v", err)
	}
	sm.clearLatestCache()

	sm.pub <- &SessionEvent{
		SessionID: latest.Id,
		Type:      SESSION_PAUSED,
	}
	requestStateReport()
	fmt.Printf("Paused session %d\n", latest.Id)

	return latest, nil
}

// EndSession will end a session if one is in progress
func (sm *SessionManager) EndSession() (*Session, error) {
	latest, err := sm.GetLatestSession()
	if err != nil {
		return nil, fmt.Errorf("failed getting current session: %v", err)
	}
	if latest == nil || latest.Complete {
		return nil, fmt.Errorf("no session in progress")
	}

	latest.Complete = true
	latest.Paused = false
	latest, err = sm.s.updateSession(latest)
	if err != nil {
		return nil, fmt.Errorf("failed to update session: %v", err)
	}
	sm.clearLatestCache()

	sm.pub <- &SessionEvent{
		SessionID: latest.Id,
		Type:      SESSION_ENDED,
	}
	fmt.Printf("Ended session %d\n", latest.Id)

	// end of session emails
	if latest.Production {
		requestCleaning(latest)
		requestPieceSelection(latest)
	}

	return latest, nil
}

// GetLatestSession returns nil, nil if there are no sessions yet
func (sm *SessionManager) GetLatestSession() (*Session, error) {
	sm.lock.Lock()
	defer sm.lock.Unlock()

	if sm.latestSessionCache != nil {
		return sm.latestSessionCache.copy(), nil
	}

	latest, err := sm.s.getLatest()
	if err != nil {
		return nil, err
	}
	sm.latestSessionCache = latest
	return latest.copy(), nil
}

// GetLatestProductionSession returns nil, nil if there are no sessions yet
func (sm *SessionManager) GetLatestProdutionSession() (*Session, error) {
	sm.lock.Lock()
	defer sm.lock.Unlock()

	latest, err := sm.s.getLatestProduction()
	if err != nil {
		return nil, err
	}
	return latest, nil
}

func (sm *SessionManager) clearLatestCache() {
	sm.lock.Lock()
	defer sm.lock.Unlock()

	sm.latestSessionCache = nil
}

func requestStateReport() {
	err := mqtt.Publish(topics_firmware.TOPIC_STATE_REPORT_REQUEST, "")
	if err != nil {
		fmt.Printf("failed to request state report from firmware: %v\n", err)
	}
}
