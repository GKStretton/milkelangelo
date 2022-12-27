package livecapture

import (
	"flag"
	"fmt"

	"github.com/gkstretton/dark/services/goo/session"
)

var (
	captureInterval = flag.Int("captureInterval", 10, "how many seconds between image captures")
	rec             *recorder
)

type recorder struct {
	sm            *session.SessionManager
	isRecording   bool
	stopRecording chan bool
}

func Run(sm *session.SessionManager) {
	rec = &recorder{
		sm:            sm,
		isRecording:   false,
		stopRecording: make(chan bool),
	}
	go rec.run()
}

func (r *recorder) run() {
	r.evaluateAction()

	// special handler for the crop config preview capture
	registerDslrPreviewHandler()

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
	if r.isRecording && (latestSession.Complete || latestSession.Paused) {
		// Stop recording
		r.stopRecording <- true
	}

	// If we're not recording but should be
	if !r.isRecording && (!latestSession.Complete && !latestSession.Paused) {
		// Start recording
		go r.record(latestSession.Id)
	}
}
