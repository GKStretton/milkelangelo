package app

import (
	"sync"

	"github.com/gkstretton/study-of-light/twitch-ebs/gooapi"
	"github.com/op/go-logging"
)

var l = logging.MustGetLogger("app")

type App struct {
	goo       gooapi.GooApi
	twitchAPI TwitchAPI

	lock  sync.Mutex
	state EBSState
}

func NewApp(goo gooapi.GooApi, twitchAPI TwitchAPI) *App {
	return &App{
		goo:       goo,
		twitchAPI: twitchAPI,
	}
}

type EBSState struct {
	GooState *gooapi.GooStateUpdate
}

func (a *App) Start() {
	a.goo.SetStateUpdateCallback(a.gooStateCallback)

	a.regularBroadcast()
}

func (a *App) gooStateCallback(state gooapi.GooStateUpdate) {
	a.state.GooState = &state

	l.Debugf("callback received state update: %+v", state)
}
