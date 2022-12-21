package livecapture

import (
	"flag"

	"github.com/gkstretton/dark/services/goo/session"
)

var (
	captureInterval = flag.Int("captureInterval", 10, "how many seconds between image captures")
	rec             *recorder
)

type recorder struct {
	sm          *session.SessionManager
	isRecording bool
	stop        chan bool
}

func Run(sm *session.SessionManager) {
	rec = &recorder{
		sm:          sm,
		isRecording: false,
		stop:        make(chan bool),
	}
	go rec.run()
}

func (r *recorder) run() {
	r.evaluateAction()

	// Listen for ongoing begin/end session events
	ch := r.sm.SubscribeToEvents()
	for {
		<-ch
		r.evaluateAction()
	}
}

func (r *recorder) evaluateAction() {
	session, _ := r.sm.GetCurrentSession()

	// If no active session and we're recording
	if session == nil && r.isRecording {
		// Stop recording
		r.stop <- true
	}

	// If there is an active session and we're not recording
	if session != nil && !r.isRecording {
		// Start recording
		go r.record()
	}
}
