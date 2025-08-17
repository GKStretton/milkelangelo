package app

import (
	"sync"
	"time"

	"github.com/gkstretton/study-of-light/twitch-ebs/entities"
	"github.com/gkstretton/study-of-light/twitch-ebs/gooapi"
	"github.com/op/go-logging"
)

var l = logging.MustGetLogger("app")

type App struct {
	goo       gooapi.GooApi
	twitchAPI TwitchAPI

	lock     sync.Mutex
	GooState *gooapi.GooStateUpdate

	ConnectedUser *entities.User
	expiryTimer   *time.Timer
}

func NewApp(goo gooapi.GooApi, twitchAPI TwitchAPI) *App {
	return &App{
		goo:       goo,
		twitchAPI: twitchAPI,
	}
}

func (a *App) Start() {
	a.goo.SetStateUpdateCallback(a.gooStateCallback)

	go a.regularTwitchStateBroadcast()
	go a.regularGooStateUpdate()
}

func (a *App) gooStateCallback(state gooapi.GooStateUpdate) {
	a.GooState = &state

	l.Debugf("callback received state update: %+v", state)
}
