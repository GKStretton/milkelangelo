package session

import (
	"fmt"
)

// memoryStorage is an in-memory, non-persistent store for testing
type memoryStorage struct {
	memoryStore map[ID]*Session
}

func (s *memoryStorage) createSession(session *Session) (*Session, error) {
	session.Id = s.getMaxID() + 1
	s.memoryStore[session.Id] = session
	return session, nil
}

func (s *memoryStorage) readSession(id ID) (*Session, error) {
	session, ok := s.memoryStore[id]
	if !ok {
		return nil, fmt.Errorf("session %d not found", id)
	}
	return session, nil
}

func (s *memoryStorage) updateSession(session *Session) (*Session, error) {
	_, ok := s.memoryStore[session.Id]
	if !ok {
		return nil, fmt.Errorf("cannot update because session %d not found", session.Id)
	}

	s.memoryStore[session.Id] = session

	return session, nil
}

func (s *memoryStorage) deleteSession(id ID) error {
	_, ok := s.memoryStore[id]
	if !ok {
		return fmt.Errorf("cannot delete because session %d not found", id)
	}

	delete(s.memoryStore, id)

	return nil
}

func (s *memoryStorage) matchSession(matcher *SessionMatcher) ([]*Session, error) {
	var matches []*Session
	for _, session := range s.memoryStore {
		// if a matcher field doesn't match, skip this one
		if matcher.Id != nil && *matcher.Id != session.Id {
			continue
		}
		if matcher.Complete != nil && *matcher.Complete != session.Complete {
			continue
		}
		if matcher.Production != nil && *matcher.Production != session.Production {
			continue
		}
		if matcher.ProductionId != nil && *matcher.ProductionId != session.ProductionId {
			continue
		}

		// add to match because matcher passed
		matches = append(matches, session)
	}

	return matches, nil
}

func (s *memoryStorage) getLatest() (*Session, error) {
	// if there's none
	if len(s.memoryStore) == 0 {
		return nil, nil
	}
	id := s.getMaxID()
	return s.readSession(id)
}

func (s *memoryStorage) getLatestProduction() (*Session, error) {
	var idOfMax ID
	var maxProductionID ID
	for id, s := range s.memoryStore {
		if s.ProductionId > maxProductionID {
			maxProductionID = s.ProductionId
			idOfMax = id
		}
	}

	// nil if there's none
	if maxProductionID == 0 {
		return nil, nil
	}

	return s.readSession(idOfMax)
}

func (s *memoryStorage) getMaxID() ID {
	var max ID
	for id := range s.memoryStore {
		if id > max {
			max = id
		}
	}
	return max
}
