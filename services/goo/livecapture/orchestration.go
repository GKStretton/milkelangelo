package livecapture

import (
	"flag"
	"fmt"
	"sync"

	"github.com/gkstretton/dark/services/goo/session"
)

var (
	captureInterval = flag.Int("captureInterval", 60, "how many seconds between image captures")
	rec             *recorder
)

type recorder struct {
	sm            *session.SessionManager
	recording     bool
	stopRecording chan bool
	mutex         *sync.RWMutex
}

func Run(sm *session.SessionManager) {
	rec = &recorder{
		sm:            sm,
		recording:     false,
		stopRecording: make(chan bool),
		mutex:         &sync.RWMutex{},
	}
	go rec.run()
}

func (r *recorder) isRecording() bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.recording
}

func (r *recorder) setIsRecording(b bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.recording = b
}

func (r *recorder) run() {
	// special handler for the crop config preview capture
	registerDslrPreviewHandler()

	r.evaluateAction()

	// Listen for ongoing begin/end session events
	ch := r.sm.SubscribeToEvents()
	for {
		<-ch
		r.evaluateAction()
	}
}

func (r *recorder) evaluateAction() {
	latestSession, err := r.sm.GetLatestSession()
	if err != nil {
		fmt.Printf("failed to GetLatestSession in livecapture: %v\n", err)
		return
	}
	if latestSession == nil {
		// no session
		return
	}

	// If we're recording but shouldn't be
	if r.isRecording() && (latestSession.Complete || latestSession.Paused) {
		// Stop recording
		r.stopRecording <- true
	}

	// If we're not recording but should be
	if !r.isRecording() && (!latestSession.Complete && !latestSession.Paused) {
		// Start recording
		go r.record(latestSession.Id)
	}
}
