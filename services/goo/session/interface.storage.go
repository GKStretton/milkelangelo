package session

type storage interface {
	createSession(session *Session) (*Session, error)
	readSession(id ID) (*Session, error)
	updateSession(session *Session) (*Session, error)
	deleteSession(id ID) error
	// matchSession returns a slice of sessions where all non-nil fields match
	matchSession(match *Session) ([]*Session, error)
}

func newStorage() storage {
	return &memoryStorage{}
}
