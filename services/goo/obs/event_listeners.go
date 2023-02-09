package obs

import (
	"fmt"

	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/dark/services/goo/events"
	"github.com/gkstretton/dark/services/goo/session"
)

func sessionListener(sm *session.SessionManager) {
	sessionChan := sm.SubscribeToEvents()
	stateReportChan := events.Subscribe()
	for {
		select {
		case <-sessionChan:
			handleSessionEvent(sm)
		case sr := <-stateReportChan:
			handleStateReport(sr)
		}
	}
}

func handleStateReport(sr *machinepb.StateReport) {
	var scene string

	if sr.Status == machinepb.Status_IDLE_MOVING ||
		sr.Status == machinepb.Status_IDLE_STATIONARY {
		scene = SCENE_IDLE
	} else {
		scene = SCENE_LIVE
	}

	err := setScene(scene)
	if err != nil {
		fmt.Printf("error setting scene in session listener: %v\n", err)
	}
}

func handleSessionEvent(e *session.SessionManager) {
	var scene string

	s, err := e.GetLatestSession()
	if err != nil {
		fmt.Printf("failed to get latest session in handleSessionEvent: %v\n", err)
		scene = SCENE_FALLBACK
	} else if s == nil {
		scene = SCENE_FALLBACK
	} else if s.Complete {
		scene = SCENE_FALLBACK
	} else if s.Paused {
		scene = SCENE_PAUSED
	} else {
		scene = SCENE_IDLE
	}

	err = setScene(scene)
	if err != nil {
		fmt.Printf("error setting scene in session listener: %v\n", err)
	}
}
