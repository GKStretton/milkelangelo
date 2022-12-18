package session

import "fmt"

type sqlStorage struct{}

func (s *sqlStorage) createSession(session *Session) (*Session, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *sqlStorage) readSession(id ID) (*Session, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *sqlStorage) updateSession(session *Session) (*Session, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *sqlStorage) deleteSession(id ID) error {
	return fmt.Errorf("not implemented")
}

func (s *sqlStorage) matchSession(match *Session) ([]*Session, error) {
	return nil, fmt.Errorf("not implemented")
}
