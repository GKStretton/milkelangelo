package livechat

import (
	"fmt"

	"github.com/gempir/go-twitch-irc/v4"
)

const channelName = "studyoflight"

// todo: this
// a twitch chat <> mqtt bridge

func Start() {
	client := twitch.NewClient(channelName, "oauth:br8dsvu2sk6qwlgvu78yb8hzjmaahb")
	// https://dev.twitch.tv/docs/authentication/refresh-tokens/
	// client.SetIRCToken()

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		for _, e := range message.Emotes {
			fmt.Printf("emote %s: https://static-cdn.jtvnw.net/emoticons/v1/%s/3.0\n", e.Name, e.ID)
		}
	})

	client.Join(channelName)

	err := client.Connect()
	if err != nil {
		panic(err)
	}
}
