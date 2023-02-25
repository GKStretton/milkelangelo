package session

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSubscribeToEvents(t *testing.T) {
	sm := NewSessionManager(true)

	// subscribe twice.
	wg := sync.WaitGroup{}

	wg.Add(1)
	var chSessionID ID
	go func() {
		defer wg.Done()
		ch := sm.SubscribeToEvents()

		t.Log("waiting for started event")
		e := <-ch
		assert.Equal(t, SESSION_STARTED, e.Type)
		chSessionID = e.SessionID

		t.Log("waiting for paused event")
		e = <-ch
		assert.Equal(t, SESSION_PAUSED, e.Type)
		assert.Equal(t, chSessionID, e.SessionID)

		t.Log("waiting for resumed event")
		e = <-ch
		assert.Equal(t, SESSION_RESUMED, e.Type)
		assert.Equal(t, chSessionID, e.SessionID)

		t.Log("waiting for ended event")
		e = <-ch
		assert.Equal(t, SESSION_ENDED, e.Type)
		assert.Equal(t, chSessionID, e.SessionID)
	}()

	sl := time.Millisecond * time.Duration(100)
	time.Sleep(sl)

	// Begin session
	startedSession, _ := sm.BeginSession(false)
	time.Sleep(sl)

	sm.PauseSession()
	time.Sleep(sl)

	sm.ResumeSession()
	time.Sleep(sl)

	sm.EndSession()
	wg.Wait()
	assert.Equal(t, startedSession.Id, chSessionID)
}

// TestSessionBasics tests BeginSession, EndSession, and GetCurrentSession
func TestSessionBasics(t *testing.T) {
	sm := NewSessionManager(true)

	// This tests BeginSession with reference to GetCurrentSession

	// No session in progress
	session, err := sm.GetLatestSession()
	assert.Nil(t, session)
	assert.NoError(t, err)

	// Begin

	startedSession, err := sm.BeginSession(false)
	assert.NoError(t, err)
	assert.Equal(t, false, startedSession.Complete)

	session, err = sm.GetLatestSession()
	assert.Equal(t, startedSession, session)
	assert.NoError(t, err)

	// Pause

	pausedSession, err := sm.PauseSession()
	assert.NoError(t, err)
	assert.True(t, pausedSession.Paused)
	assert.False(t, pausedSession.Complete)

	// Resume
	resumedSession, err := sm.ResumeSession()
	assert.NoError(t, err)
	assert.False(t, resumedSession.Paused)
	assert.False(t, resumedSession.Complete)

	// End

	endedSession, err := sm.EndSession()
	assert.NoError(t, err)
	assert.True(t, endedSession.Complete)

	// No session in progress
	session, err = sm.GetLatestSession()
	assert.NoError(t, err)
	assert.True(t, session.Complete)
}
