package app

import (
	"time"

	"github.com/gkstretton/study-of-light/twitch-ebs/gooapi"
)

func (a *App) regularStateUpdate() {
	next := time.After(time.Second)
	for {
		<-next
		next = time.After(time.Second)

		a.reportEbsState()
	}
}

func (a *App) reportEbsState() {
	a.lock.Lock()
	defer a.lock.Unlock()

	err := a.goo.ReportEbsState(gooapi.EbsStateReport{
		ConnectedUser: a.ConnectedUser,
	})
	if err != nil {
		l.Errorf("failed to report ebs state to goo: %v", err)
	}
}
