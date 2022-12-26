package session

import (
	"fmt"
	"path/filepath"

	"github.com/gkstretton/dark/services/goo/filesystem"
)

type yamlStorage struct{}

func (s *yamlStorage) createSession(session *Session) (*Session, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *yamlStorage) readSession(id ID) (*Session, error) {

	return nil, fmt.Errorf("not implemented")
}

func (s *yamlStorage) updateSession(session *Session) (*Session, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *yamlStorage) deleteSession(id ID) error {
	return fmt.Errorf("not implemented")
}

func (s *yamlStorage) matchSession(matcher *SessionMatcher) ([]*Session, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *yamlStorage) getLatest() (*Session, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *yamlStorage) getSessionYamlPath(id ID) string {
	return filepath.Join(filesystem.GetMetadataDir(uint64(id)), "session.yml")
}
