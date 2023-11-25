package twitch

import (
	"fmt"
	"os"

	"github.com/nicklaw5/helix/v2"
)

type Colour string

const (
	COLOUR_PRIMARY Colour = "primary"
	COLOUR_BLUE    Colour = "blue"
	COLOUR_GREEN   Colour = "green"
	COLOUR_ORANGE  Colour = "orange"
	COLOUR_PURPLE  Colour = "purple"
)

func (c *TwitchApi) Announce(msg string, colour Colour) {
	resp, err := c.helixClient.SendChatAnnouncement(&helix.SendChatAnnouncementParams{
		BroadcasterID: channelId,
		ModeratorID:   channelId,
		Message:       msg,
		Color:         string(colour),
	})
	if err != nil || resp.Error != "" {
		fmt.Printf("failed to announce '%s': %v %v\n", msg, err, resp)
	}
}

func (c *TwitchApi) Say(msg string) {
	c.ircClient.Say(channelName, msg)
}

func (c *TwitchApi) Reply(msgId string, msg string) {
	c.ircClient.Reply(channelName, msgId, msg)
}

func (c *TwitchApi) SendExtensionMessage(payload string) {
	r, err := c.helixClient.SendExtensionPubSubMessage(&helix.ExtensionSendPubSubMessageParams{
		BroadcasterID: channelId,
		Message:       payload,
		Target:        []helix.ExtensionPubSubPublishType{helix.ExtensionPubSubBroadcastPublish},
	})
	if err != nil {
		fmt.Printf("failed to send pubsub message: %v\n", err)
		return
	}
	if r.Error != "" {
		fmt.Printf("failed to send pubsub message: %v\n", r)
		return
	}

}

func init() {
	api := Start()
	api.SendExtensionMessage("test msg")
	os.Exit(0)
}
