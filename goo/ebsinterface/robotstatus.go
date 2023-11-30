package ebsinterface

import "github.com/gkstretton/dark/services/goo/events"

type robotStatus struct {
	Status string
}

func (e *ExtensionSession) regularRobotStatusUpdate() {
	c := events.Subscribe()
	for {
		select {
		case <-e.exitCh:
			l.Println("exiting regular robot status update loop")
			return
		case sr := <-c:
			e.updateRobotStatus(&robotStatus{
				Status: sr.Status.String(),
			})
		}
	}
}
