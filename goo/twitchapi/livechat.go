package twitchapi

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gempir/go-twitch-irc/v4"
	"github.com/gkstretton/dark/services/goo/keyvalue"
	"github.com/nicklaw5/helix/v2"
)

const channelName = "studyoflight"
const channelId = "807784320"

type TwitchApi struct {
	helixClient *helix.Client
	ircClient   *twitch.Client

	subs []chan *Message
	lock sync.Mutex
}

func Start() *TwitchApi {
	// auth
	clientID := string(keyvalue.Get("TWITCH_CLIENT_ID"))
	clientSecret := string(keyvalue.Get("TWITCH_CLIENT_SECRET"))

	if clientID == "" || clientSecret == "" {
		fmt.Println("error: twitch client id or secret not set")
		return nil
	}

	// Helix api
	helixClient, err := helix.NewClient(&helix.Options{
		ClientID:        clientID,
		ClientSecret:    clientSecret,
		RefreshToken:    getRefreshToken(),
		UserAccessToken: "x", // will auto refresh
		HTTPClient: &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:    10,
				IdleConnTimeout: 30 * time.Second,
			},
			Timeout: 10 * time.Second,
		},
	})
	if err != nil {
		fmt.Printf("failed to create twitch helix client: %v\n", err)
		return nil
	}

	// IRC api
	ircClient := twitch.NewClient(channelName, "oauth:x") // will be refreshed by helix client

	helixClient.OnUserAccessTokenRefreshed(func(newAccessToken, newRefreshToken string) {
		// Update the client's token
		ircClient.SetIRCToken("oauth:" + newAccessToken)
		// Update the stored refresh token with the new one
		setRefreshToken(newRefreshToken)
		fmt.Println("twitch helix client refreshed access token")
	})

	// will trigger auth refresh
	go func() {
		_, err = helixClient.SendChatAnnouncement(&helix.SendChatAnnouncementParams{
			BroadcasterID: channelId,
			ModeratorID:   channelId,
			Message:       "Milkelangelo backend is now running",
			// value must be one of "", "primary", "purple", "blue", "green", "orange"
			Color: "primary",
		})
		if err != nil {
			fmt.Printf("failed to send initial helix broadcast: %v\n", err)
		}
		fmt.Println("sent twitch startup announcement.")
	}()

	api := &TwitchApi{
		helixClient: helixClient,
		ircClient:   ircClient,
		subs:        []chan *Message{},
	}

	// Set up your client, handlers, etc. here
	api.setupHandlers()

	// listen to irc
	go func() {
		ircClient.Join(channelName)
		fmt.Println("connecting to twitch irc")
		err = ircClient.Connect()
		if err != nil {
			fmt.Printf("failed to listen to twitch irc: %v\n", err)
		}
	}()

	return api
}

func getRefreshToken() string {
	return string(keyvalue.Get("TWITCH_REFRESH_TOKEN"))
}

func setRefreshToken(token string) {
	keyvalue.Set("TWITCH_REFRESH_TOKEN", []byte(token))
}
