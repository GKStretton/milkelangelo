package ebsinterface

import "github.com/gkstretton/dark/services/goo/events"

type robotStatus struct {
	Status string
}

func (e *extensionSession) regularRobotStatusUpdate() {
	c := events.Subscribe()
	for {
		select {
		case <-e.exitCh:
			l.Println("exiting regular robot status update loop")
			return
		case sr := <-c:
			if sr == nil {
				e.updateRobotStatus(nil)
				continue
			}
			e.updateRobotStatus(&robotStatus{
				Status: sr.Status.String(),
			})
		}
	}
}
