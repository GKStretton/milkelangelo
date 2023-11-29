package ebsinterface

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gkstretton/dark/services/goo/keyvalue"
)

// consists of latest state report and current vote status
type BroadcastData = string

func (e *extensionSession) BroadcastData(data BroadcastData) error {
	type payload struct {
		message        string
		broadcaster_id string
		target         []string
	}
	pl := &payload{
		message:        data,
		broadcaster_id: channelId,
		target:         []string{"broadcast"},
	}

	jsonData, err := json.Marshal(pl)
	if err != nil {
		return err
	}
	url := "https://api.twitch.tv/helix/extensions/pubsub"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	clientID := string(keyvalue.Get("TWITCH_EXTENSION_CLIENT_ID"))

	req.Header.Set("Authorization", "Bearer "+e.broadcastToken)
	req.Header.Set("Client-Id", clientID)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
