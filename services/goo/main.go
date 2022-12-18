package main

import (
	"github.com/gkstretton/dark/services/goo/mqtt"
	"github.com/gkstretton/dark/services/goo/session"
)

func main() {
	mqtt.Start()

	sm := session.NewSessionManager()
	<-sm.GetEventsChan()
}
