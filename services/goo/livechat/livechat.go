package livechat

import (
	"fmt"
	"time"

	"github.com/gempir/go-twitch-irc/v4"
	"github.com/gkstretton/dark/services/goo/keyvalue"
)

const channelName = "studyoflight"

func Start() {
	clientID := string(keyvalue.Get("TWITCH_CLIENT_ID"))
	clientSecret := string(keyvalue.Get("TWITCH_CLIENT_SECRET"))

	// Obtain a new access token using the refresh token
	tokenResponse, err := RefreshAccessToken(clientID, clientSecret, getRefreshToken())
	if err != nil {
		fmt.Println(err)
		return
	}
	setRefreshToken(tokenResponse.RefreshToken)

	client := twitch.NewClient(channelName, "oauth:"+tokenResponse.AccessToken)

	go refreshTokenPeriodically(client, clientID, clientSecret)

	// Set up your client, handlers, etc. here
	setup(client)

	client.Join(channelName)

	err = client.Connect()
	if err != nil {
		panic(err)
	}
}

func refreshTokenPeriodically(client *twitch.Client, clientID, clientSecret string) {
	interval := 3*time.Hour + 30*time.Minute
	next := time.After(interval) // Set the interval to 3.5 hours
	for {
		<-next

		// Call the function to refresh the token here
		tokenResponse, err := RefreshAccessToken(clientID, clientSecret, getRefreshToken())
		if err != nil {
			fmt.Printf("Could not refresh the twitch access token: %v", err)
			next = time.After(time.Minute * time.Duration(5))
			continue
		}
		next = time.After(interval)

		// Update the client's token
		client.SetIRCToken("oauth:" + tokenResponse.AccessToken)
		// Update the stored refresh token with the new one
		setRefreshToken(tokenResponse.RefreshToken)
	}
}

func setup(client *twitch.Client) {
	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		for _, e := range message.Emotes {
			fmt.Printf("emote %s: https://static-cdn.jtvnw.net/emoticons/v1/%s/3.0\n", e.Name, e.ID)
		}
		fmt.Println("got message")
	})
	go func() {
		time.Sleep(time.Second * 5)
		client.Say(channelName, "goo/livechat is running")
	}()
}
