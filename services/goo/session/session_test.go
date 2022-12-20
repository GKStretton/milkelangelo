package session

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSubscribeToEvents(t *testing.T) {
	//todo: subscribe twice.
	//todo: begin session
	//todo: end session
}

func TestEndSession(t *testing.T) {
}

func TestBeginSession(t *testing.T) {
	//todo: cut down

	sm := NewSessionManager(true)

	// This tests the event listener feature

	wg := sync.WaitGroup{}

	wg.Add(1)
	var chSessionID1 ID
	go func() {
		defer wg.Done()
		ch, err := sm.SubscribeToEvents()
		assert.NoError(t, err)

		e := <-ch
		assert.Equal(t, SESSION_STARTED, e.Type)
		chSessionID1 = e.SessionID
	}()

	wg.Add(1)
	var chSessionID2 ID
	go func() {
		defer wg.Done()
		ch, err := sm.SubscribeToEvents()
		assert.NoError(t, err)

		e := <-ch
		assert.Equal(t, SESSION_STARTED, e.Type)
		chSessionID2 = e.SessionID
	}()

	// ensure they are subscribed
	time.Sleep(time.Millisecond * time.Duration(10))

	// This tests BeginSession with reference to GetCurrentSession

	// No session in progress
	session, err := sm.GetCurrentSession()
	assert.Nil(t, session)
	assert.NoError(t, err)

	startedSession, err := sm.BeginSession()
	assert.NoError(t, err)
	assert.Equal(t, false, startedSession.Complete)

	// No session in progress
	retrievedSession, err := sm.GetCurrentSession()
	assert.Equal(t, startedSession, retrievedSession)
	assert.NoError(t, err)

	wg.Wait()
	assert.Equal(t, retrievedSession.Id, chSessionID1)
	assert.Equal(t, retrievedSession.Id, chSessionID2)
}
