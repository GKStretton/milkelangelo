package app

import (
	"encoding/json"
	"time"

	"github.com/gkstretton/study-of-light/twitch-ebs/entities"
	"github.com/gkstretton/study-of-light/twitch-ebs/gooapi"
)

// broadcasts the BroadcastData cache once per second
func (a *App) regularBroadcast() {
	// get marshaled data, protected by lock
	d := func() ([]byte, error) {
		a.lock.Lock()
		defer a.lock.Unlock()

		state := struct {
			GooState               *gooapi.GooStateUpdate
			ConnectedUser          *entities.User
			ConnectedUserTimestamp time.Time
		}{
			GooState:               a.GooState,
			ConnectedUser:          a.ConnectedUser,
			ConnectedUserTimestamp: a.ConnectedUserTimestamp,
		}

		jsonData, err := json.Marshal(state)
		if err != nil {
			return nil, err
		}
		return jsonData, nil
	}

	send := func() {
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

	next := time.After(0)
	for {
		<-next
		next = time.After(time.Millisecond * 1000)
		send()
	}
}
