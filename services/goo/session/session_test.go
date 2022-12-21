package session

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubscribeToEvents(t *testing.T) {
	sm := NewSessionManager(true)

	// subscribe twice.
	wg := sync.WaitGroup{}

	wg.Add(1)
	var chSessionIDStarted ID
	var chSessionIDEnded ID
	go func() {
		defer wg.Done()
		ch := sm.SubscribeToEvents()

		e := <-ch
		assert.Equal(t, SESSION_STARTED, e.Type)
		chSessionIDStarted = e.SessionID

		e = <-ch
		assert.Equal(t, SESSION_ENDED, e.Type)
		chSessionIDEnded = e.SessionID
	}()

	// Begin session
	startedSession, err := sm.BeginSession()
	assert.NoError(t, err)
	assert.Equal(t, false, startedSession.Complete)

	// End session
	endedSession, err := sm.EndSession()
	assert.NoError(t, err)
	assert.Equal(t, true, endedSession.Complete)

	wg.Wait()
	assert.Equal(t, startedSession.Id, endedSession.Id)
	assert.Equal(t, startedSession.Id, chSessionIDStarted)
	assert.Equal(t, endedSession.Id, chSessionIDEnded)
}

// TestSessionBasics tests BeginSession, EndSession, and GetCurrentSession
func TestSessionBasics(t *testing.T) {
	sm := NewSessionManager(true)

	// This tests BeginSession with reference to GetCurrentSession

	// No session in progress
	session, err := sm.GetCurrentSession()
	assert.Nil(t, session)
	assert.NoError(t, err)

	startedSession, err := sm.BeginSession()
	assert.NoError(t, err)
	assert.Equal(t, false, startedSession.Complete)

	session, err = sm.GetCurrentSession()
	assert.Equal(t, startedSession, session)
	assert.NoError(t, err)

	endedSession, err := sm.EndSession()
	assert.NoError(t, err)
	assert.Equal(t, true, endedSession.Complete)

	// No session in progress
	session, err = sm.GetCurrentSession()
	assert.Nil(t, session)
	assert.NoError(t, err)
}
