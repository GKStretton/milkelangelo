package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/gkstretton/dark/services/goo/app"
	"github.com/gkstretton/dark/services/goo/config"
	"github.com/gkstretton/dark/services/goo/contentscheduler"
	"github.com/gkstretton/dark/services/goo/control"
	"github.com/gkstretton/dark/services/goo/email"
	"github.com/gkstretton/dark/services/goo/events"
	"github.com/gkstretton/dark/services/goo/filesystem"
	"github.com/gkstretton/dark/services/goo/keyvalue"
	"github.com/gkstretton/dark/services/goo/livecapture"
	"github.com/gkstretton/dark/services/goo/mqtt"
	"github.com/gkstretton/dark/services/goo/obs"
	"github.com/gkstretton/dark/services/goo/publicapi"
	"github.com/gkstretton/dark/services/goo/server"
	"github.com/gkstretton/dark/services/goo/session"
	"github.com/gkstretton/dark/services/goo/socialmedia"
	"github.com/gkstretton/dark/services/goo/twitchapi"
	"github.com/gkstretton/dark/services/goo/vialprofiles"
	"github.com/joho/godotenv"
)

var (
	test                      = flag.Bool("test", false, "if true, just run test code")
	refreshYoutubeCredentials = flag.Bool("yt", false, "if true, refresh youtube credentials")
)

func main() {
	flag.Parse()

	godotenv.Load()

	if *refreshYoutubeCredentials {
		socialmedia.RefreshYoutubeCreds()
		return
	}

	if *test {
		runAdHocTests()
		return
	}

	filesystem.AssertBasePaths()

	mqtt.Start(config.BrokerHost())
	keyvalue.Start()
	email.Start()
	server.Start()

	sm := session.NewSessionManager(false)
	twitchApi := twitchapi.Start()

	api := publicapi.New()

	events.Start(sm, api)
	control.Start(api)
	livecapture.Start(sm)
	obs.Start(config.BrokerHost(), sm)
	vialprofiles.Start(sm, api)
	contentscheduler.Start(sm)

	app.Start(sm, twitchApi)

	if config.EnablePublicApi() {
		err := api.Start(publicapi.Options{
			Addr:        config.PublicApiAddr(),
			UiDir:       config.PublicUiDir(),
			MediamtxURL: config.MediamtxURL(),

			CloudflareTurnKeyID:    config.CloudflareTurnKeyID(),
			CloudflareTurnAPIToken: config.CloudflareTurnAPIToken(),
		})
		if err != nil {
			// don't take down the rest of goo for a misconfigured public api
			fmt.Printf("failed to start public api: %v\n", err)
		}
	}

	// Block to prevent early quit
	fmt.Println("finished init, main loop sleeping.")
	for {
		time.Sleep(time.Millisecond * time.Duration(100))
	}
}
