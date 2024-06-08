package obs

import (
	"fmt"

	"github.com/gkstretton/dark/services/goo/session"
)

func sessionListener(sm *session.SessionManager) {
	sessionChan := sm.SubscribeToEvents()

	for {
		<-sessionChan
		handleSessionEvent(sm)
	}
}

func handleSessionEvent(e *session.SessionManager) {
	var scene string

	s, err := e.GetLatestSession()
	if err != nil {
		fmt.Printf("failed to get latest session in handleSessionEvent: %v\n", err)
		scene = SCENE_ERROR
	} else if s == nil {
		scene = SCENE_COMPLETE
	} else if s.Complete {
		scene = SCENE_COMPLETE
	} else if s.Paused {
		scene = SCENE_PAUSED
	} else {
		scene = SCENE_LIVE
	}

	err = setScene(scene)
	if err != nil {
		fmt.Printf("error setting scene in session listener: %v\n", err)
	}

	if s != nil {
		if s.Production {
			setSessionNumber(int(s.ProductionId), true)
		} else {
			setSessionNumber(int(s.Id), false)
		}
	}

}
