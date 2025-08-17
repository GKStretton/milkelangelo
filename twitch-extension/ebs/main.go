package main

import (
	"flag"

	"github.com/gkstretton/study-of-light/twitch-ebs/app"
	"github.com/gkstretton/study-of-light/twitch-ebs/config"
	"github.com/gkstretton/study-of-light/twitch-ebs/gooapi"
	"github.com/gkstretton/study-of-light/twitch-ebs/server"
	"github.com/gkstretton/study-of-light/twitch-ebs/twitchapi"
	"github.com/op/go-logging"
)

var (
	addr              = flag.String("addr", ":8789", "address the public server should listen on")
	internalAddr      = flag.String("internalAddr", ":8788", "address to listen for internal (goo) requests on")
	channelID         = flag.String("channelID", "807784320", "twitch channel id")
	extensionClientID = flag.String("extensionClientID", "ihiyqlxtem517wq76f4hn8pvo9is30", "twitch extension client id")

	l = logging.MustGetLogger("ebs")
)

func main() {
	flag.Parse()

	goo, err := gooapi.NewConnectedGooApi(config.SharedSecretGoo(), *internalAddr)
	if err != nil {
		l.Fatalf("failed to create goo api: %s\n", err)
	}

	twitchAPI, err := twitchapi.NewConnectedTwitchAPI(config.SharedSecretTwitch(), *channelID, *extensionClientID)
	if err != nil {
		l.Fatalf("failed to create twitch api: %s\n", err)
	}

	// listen for internal (goo) connections
	go goo.Start()

	app := app.NewApp(goo, twitchAPI)
	app.Start()

	s, err := server.NewServer(*addr, config.SharedSecretTwitch(), goo, app)
	if err != nil {
		l.Fatalf("failed to create server: %s\n", err)
	}

	// listen to twitch clients
	s.Run()
}
