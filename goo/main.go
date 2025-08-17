package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/gkstretton/dark/services/goo/actor"
	"github.com/gkstretton/dark/services/goo/app"
	"github.com/gkstretton/dark/services/goo/config"
	"github.com/gkstretton/dark/services/goo/contentscheduler"
	"github.com/gkstretton/dark/services/goo/ebsinterface"
	"github.com/gkstretton/dark/services/goo/email"
	"github.com/gkstretton/dark/services/goo/events"
	"github.com/gkstretton/dark/services/goo/filesystem"
	"github.com/gkstretton/dark/services/goo/keyvalue"
	"github.com/gkstretton/dark/services/goo/livecapture"
	"github.com/gkstretton/dark/services/goo/mqtt"
	"github.com/gkstretton/dark/services/goo/obs"
	"github.com/gkstretton/dark/services/goo/server"
	"github.com/gkstretton/dark/services/goo/session"
	"github.com/gkstretton/dark/services/goo/socialmedia"
	"github.com/gkstretton/dark/services/goo/twitchapi"
	"github.com/gkstretton/dark/services/goo/util"
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

	var ebsApi ebsinterface.EbsApi

	if util.EnvBool("ENABLE_EBS") {
		host := keyvalue.GetString("EBS_HOST")
		if host == "" {
			panic("EBS_HOST not set")
		}

		var err error
		ebsApi, err = ebsinterface.NewExtensionSession("http://" + host + ":8788")
		if err != nil {
			panic("failed to init ebs: " + err.Error())
		}
	}

	sm := session.NewSessionManager(false)
	twitchApi := twitchapi.Start()

	actor.Setup(sm, ebsApi)
	events.Start(sm, ebsApi)
	livecapture.Start(sm)
	obs.Start(config.BrokerHost(), sm)
	vialprofiles.Start(sm, ebsApi)
	contentscheduler.Start(sm)

	app.Start(sm, twitchApi, ebsApi)

	// Block to prevent early quit
	fmt.Println("finished init, main loop sleeping.")
	for {
		time.Sleep(time.Millisecond * time.Duration(100))
	}
}
