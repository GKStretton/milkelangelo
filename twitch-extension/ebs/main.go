package main

import (
	"flag"

	"github.com/gkstretton/study-of-light/twitch-ebs/gooapi"
	"github.com/gkstretton/study-of-light/twitch-ebs/server"
	"github.com/gkstretton/study-of-light/twitch-ebs/twitchapi"
	"github.com/op/go-logging"
)

var (
	addr               = flag.String("addr", ":8789", "address the public server should listen on")
	internalAddr       = flag.String("internalAddr", ":8788", "address to listen for internal (goo) requests on")
	internalSecretPath = flag.String("internalSecretPath", "/mnt/md0/light-stores/kv/EBS_INTERNAL_SECRET", "path for secret used for internal jwt verification")
	sharedSecretPath   = flag.String("sharedSecretPath", ".shared-secret", "path for shared secret used for twitch jwt verification")
	channelID          = flag.String("channelID", "807784320", "twitch channel id")
	extensionClientID  = flag.String("extensionClientID", "ihiyqlxtem517wq76f4hn8pvo9is30", "twitch extension client id")

	l = logging.MustGetLogger("ebs")
)

func main() {
	flag.Parse()

	goo, err := gooapi.NewConnectedGooApi(*internalSecretPath, *internalAddr)
	if err != nil {
		l.Fatalf("failed to create goo api: %w\n", err)
	}

	twitchapi, err := twitchapi.NewConnectedTwitchAPI(*sharedSecretPath, *channelID, *extensionClientID)
	if err != nil {
		l.Fatalf("failed to create twitch api: %w\n", err)
	}
	_ = twitchapi

	s, err := server.NewServer(*addr, *sharedSecretPath, goo)

	if err != nil {
		l.Fatalf("failed to create server: %w\n", err)
	}

	// listen for internal (goo) connections
	go goo.Start()

	// listen to twitch clients
	s.Run()
}
