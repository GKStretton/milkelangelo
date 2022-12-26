package session

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/gkstretton/dark/services/goo/filesystem"
	"gopkg.in/yaml.v3"
)

type yamlStorage struct {
	lock *sync.Mutex
}

func (s *yamlStorage) createSession(session *Session) (*Session, error) {
	// first id
	session.Id = 1

	// overwrite if others present
	latest, _ := s.getLatest()
	if latest != nil {
		session.Id = latest.Id + 1
	}

	session, err := s.updateSession(session)
	if err != nil {
		return nil, fmt.Errorf("failed to write session yml: %v", err)
	}
	// update latest
	s.lock.Lock()
	defer s.lock.Unlock()
	err = filesystem.CreateSymlink(
		s.getSessionYamlPath(session.Id),
		s.getLatestSessionYamlPath(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update latest session metadata symlink: %v", err)
	}

	return session, nil
}

func (s *yamlStorage) readSession(id ID) (*Session, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	f, err := os.Open(s.getSessionYamlPath(id))
	if err != nil {
		return nil, fmt.Errorf("failed to open session yml: %v", err)
	}

	session := &Session{}
	e := yaml.NewDecoder(f)
	err = e.Decode(session)
	if err != nil {
		return nil, fmt.Errorf("failed to decode yml to struct: %v", err)
	}

	return session, nil
}

func (s *yamlStorage) updateSession(session *Session) (*Session, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	// open file
	f, err := os.OpenFile(s.getSessionYamlPath(session.Id), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open session yml: %v", err)
	}
	defer f.Close()

	// encode
	e := yaml.NewEncoder(f)
	err = e.Encode(session)
	if err != nil {
		return nil, fmt.Errorf("error encoding yml: %v", err)
	}
	e.Close()
	fmt.Println("written session yaml")

	return session, nil
}

func (s *yamlStorage) deleteSession(id ID) error {
	return fmt.Errorf("not implemented")
}

func (s *yamlStorage) matchSession(matcher *SessionMatcher) ([]*Session, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *yamlStorage) getLatest() (*Session, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	path := s.getLatestSessionYamlPath()

	// return nil,nil if non-existence
	if _, err := os.Lstat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to stat symlink: %v", err)
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open session yml: %v", err)
	}

	session := &Session{}
	e := yaml.NewDecoder(f)
	err = e.Decode(session)
	if err != nil {
		return nil, fmt.Errorf("failed to decode yml to struct: %v", err)
	}

	return session, nil
}

func (s *yamlStorage) getSessionYamlPath(id ID) string {
	return filepath.Join(
		filesystem.GetMetadataDir(),
		strconv.Itoa(int(id))+"_session.yml",
	)
}

func (s *yamlStorage) getLatestSessionYamlPath() string {
	return filepath.Join(filesystem.GetMetadataDir(), "latest"+"_session.yml")
}
