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
	s.lock.Lock()
	defer s.lock.Unlock()

	// first id
	session.Id = 1

	// overwrite if others present
	latest, _ := getLatest_internal()
	if latest != nil {
		session.Id = latest.Id + 1
	}

	if session.Production {
		session.ProductionId = 1
		latest, _ := getLatestProduction_internal()
		if latest != nil {
			session.ProductionId = latest.ProductionId + 1
		}
	}

	err := writeSessionToYaml(session, getSessionYamlPath(session.Id))
	if err != nil {
		return nil, fmt.Errorf("failed to write session yml: %v", err)
	}
	// update latest
	err = filesystem.CreateSymlink(
		getSessionYamlPath(session.Id),
		getLatestSessionYamlPath(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update latest session metadata symlink: %v", err)
	}

	if session.Production {
		err = filesystem.CreateSymlink(
			getSessionYamlPath(session.Id),
			getLatestProductionSessionYamlPath(),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to update latest session metadata symlink: %v", err)
		}
	}

	return session, nil
}

func (s *yamlStorage) readSession(id ID) (*Session, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	f, err := os.Open(getSessionYamlPath(id))
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
	err := writeSessionToYaml(session, getSessionYamlPath(session.Id))
	if err != nil {
		return nil, err
	}

	return session, nil
}

func writeSessionToYaml(session *Session, path string) error {
	// open file
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("failed to open session yml: %v", err)
	}
	defer f.Close()

	// encode
	e := yaml.NewEncoder(f)
	err = e.Encode(session)
	if err != nil {
		return fmt.Errorf("error encoding yml: %v", err)
	}
	e.Close()
	fmt.Println("written session yaml")
	return nil
}

func (s *yamlStorage) deleteSession(id ID) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	return fmt.Errorf("not implemented")
}

func (s *yamlStorage) matchSession(matcher *SessionMatcher) ([]*Session, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	return nil, fmt.Errorf("not implemented")
}

func (s *yamlStorage) getLatest() (*Session, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	return getLatest_internal()
}

// internal one without mutex for use in synchronised functions
func getLatest_internal() (*Session, error) {
	return loadSessionFromSymlink(
		getLatestSessionYamlPath(),
	)
}

func (s *yamlStorage) getLatestProduction() (*Session, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	return getLatestProduction_internal()
}

// internal one without mutex for use in synchronised functions
func getLatestProduction_internal() (*Session, error) {
	return loadSessionFromSymlink(
		getLatestProductionSessionYamlPath(),
	)
}

func loadSessionFromSymlink(path string) (*Session, error) {
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

func getSessionYamlPath(id ID) string {
	return filepath.Join(
		filesystem.GetMetadataDir(),
		strconv.Itoa(int(id))+"_session.yml",
	)
}

func getLatestSessionYamlPath() string {
	return filepath.Join(filesystem.GetMetadataDir(), "latest_session.yml")
}

func getLatestProductionSessionYamlPath() string {
	return filepath.Join(filesystem.GetMetadataDir(), "latest_production_session.yml")
}
