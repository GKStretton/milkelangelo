package main

import (
	"bufio"
	"flag"
	"os"

	"github.com/gkstretton/dark/services/goo/livecapture"
	"github.com/gkstretton/dark/services/goo/mqtt"
	"github.com/gkstretton/dark/services/goo/session"
)

func main() {
	flag.Parse()

	mqtt.Start()

	sm := session.NewSessionManager(true)

	livecapture.Run(sm)

	// Block to prevent early quit
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
}
