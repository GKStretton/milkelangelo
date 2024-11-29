package twitchapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gkstretton/study-of-light/twitch-ebs/common"
)

type connectedTwitchAPI struct {
	sharedSecret      []byte
	channelID         string
	extensionClientID string

	broadcastToken string
}

func NewConnectedTwitchAPI(sharedSecretPath, channelID, extensionClientID string) (*connectedTwitchAPI, error) {
	sharedSecret, err := common.GetSecret(sharedSecretPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get shared secret: %w", err)
	}

	t := &connectedTwitchAPI{
		sharedSecret:      sharedSecret,
		channelID:         channelID,
		extensionClientID: extensionClientID,
	}

	//! Note token refreshing is not yet implemented!
	//! It's assumed the ebs will be turned off between sessions.
	bt, err := t.getBroadcastToken(time.Hour * 2)
	if err != nil {
		return nil, fmt.Errorf("failed to generate broadcast token: %w", err)
	}
	t.broadcastToken = bt

	return t, nil
}

// broadcastData must be called with rate limiting due to pubsub api limit.
// This is officially stated as 100 per minute, but there's a thread saying it's
// 1 regen per second with pool of 100.
// https://github.com/twitchdev/issues/issues/612
// So we should stick to 1 per second, 60 per minute.
func (t *connectedTwitchAPI) BroadcastExtensionData(jsonData []byte) error {
	type payload struct {
		Message       string   `json:"message"`
		BroadcasterID string   `json:"broadcaster_id"`
		Target        []string `json:"target"`
	}
	pl := &payload{
		Message:       string(jsonData),
		BroadcasterID: t.channelID,
		Target:        []string{"broadcast"},
	}

	jsonPl, err := json.Marshal(pl)
	if err != nil {
		return err
	}
	url := "https://api.twitch.tv/helix/extensions/pubsub"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPl))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+t.broadcastToken)
	req.Header.Set("Client-Id", t.extensionClientID)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error reading error response: %v", err)
		}
		return errors.New(string(b))
	}

	return nil
}
