package session

import "sync"

type storage interface {
	createSession(session *Session) (*Session, error)
	readSession(id ID) (*Session, error)
	updateSession(session *Session) (*Session, error)
	deleteSession(id ID) error
	// matchSession returns a slice of sessions where all non-nil fields match
	matchSession(matcher *SessionMatcher) ([]*Session, error)
	// getLatest returns the latest to be added (highest id)
	getLatest() (*Session, error)
}

func newStorage(useMemoryStorage bool) storage {
	if useMemoryStorage {
		return &memoryStorage{
			memoryStore: map[ID]*Session{},
		}
	} else {
		return &yamlStorage{
			lock: &sync.Mutex{},
		}
	}
}
