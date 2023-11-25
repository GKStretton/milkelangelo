package twitch

import (
	"fmt"

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
