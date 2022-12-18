package session

import "fmt"

// memoryStorage is an in-memory, non-persistent store for testing
type memoryStorage struct{}

func (s *memoryStorage) createSession(session *Session) (*Session, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *memoryStorage) readSession(id ID) (*Session, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *memoryStorage) updateSession(session *Session) (*Session, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *memoryStorage) deleteSession(id ID) error {
	return fmt.Errorf("not implemented")
}

func (s *memoryStorage) matchSession(match *Session) ([]*Session, error) {
	return nil, fmt.Errorf("not implemented")
}
