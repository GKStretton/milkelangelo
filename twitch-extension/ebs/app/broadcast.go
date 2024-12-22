package app

import (
	"encoding/json"
	"time"

	"github.com/gkstretton/study-of-light/twitch-ebs/entities"
	"github.com/gkstretton/study-of-light/twitch-ebs/gooapi"
)

type ebsState struct {
	GooState      *gooapi.GooStateUpdate
	ConnectedUser *entities.User
}

func (a *App) buildStateResponse() ebsState {
	return ebsState{
		GooState:      a.GooState,
		ConnectedUser: a.ConnectedUser,
	}
}

// broadcasts the BroadcastData cache once per second
func (a *App) regularBroadcast() {
	next := time.After(0)
	for {
		<-next
		next = time.After(time.Millisecond * 1000)
		a.broadcast()
	}
}

func (a *App) broadcast() {
	// get marshaled data, protected by lock
	d := func() ([]byte, error) {
		a.lock.Lock()
		defer a.lock.Unlock()

		jsonData, err := json.Marshal(a.buildStateResponse())
		if err != nil {
			return nil, err
		}
		return jsonData, nil
	}

	data, err := d()
	if err != nil {
		l.Errorf("failed to marshal broadcast data: %v\n", err)
		return
	}
	err = a.twitchAPI.BroadcastExtensionData(data)
	if err != nil {
		l.Errorf("failed to send broadcast data: %v\n", err)
		return
	}
}
