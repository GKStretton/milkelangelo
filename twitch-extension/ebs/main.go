package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gkstretton/study-of-light/twitch-ebs/gooapi"
	"github.com/gkstretton/study-of-light/twitch-ebs/server"
)

var (
	addr               = flag.String("addr", ":8789", "address the public server should listen on")
	internalAddr       = flag.String("internalAddr", ":8788", "address to listen for internal (goo) requests on")
	internalSecretPath = flag.String("internalSecretPath", ".internal-secret", "path for secret used for internal jwt verification")
	sharedSecretPath   = flag.String("sharedSecretPath", ".shared-secret", "path for shared secret used for twitch jwt verification")
)

func main() {
	flag.Parse()

	goo, err := gooapi.NewConnectedGooApi(*internalSecretPath, *internalAddr)
	if err != nil {
		fmt.Printf("failed to create goo api: %v\n", err)
		os.Exit(1)
	}

	s, err := server.NewServer(*addr, *sharedSecretPath, goo)
	if err != nil {
		fmt.Printf("failed to create server: %v\n", err)
		os.Exit(1)
	}

	// listen for internal (goo) connections
	go goo.Start()

	// listen to twitch clients
	s.Run()
}
