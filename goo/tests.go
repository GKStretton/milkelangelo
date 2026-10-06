package main

import (
	"fmt"

	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/dark/services/goo/config"
	"github.com/gkstretton/dark/services/goo/control"
	"github.com/gkstretton/dark/services/goo/email"
	"github.com/gkstretton/dark/services/goo/events"
	"github.com/gkstretton/dark/services/goo/mqtt"
	"github.com/gkstretton/dark/services/goo/publicapi"
	"github.com/gkstretton/dark/services/goo/session"
	"github.com/gkstretton/dark/services/goo/vialprofiles"
)

// tests for human verification during development
func runAdHocTests() {
	testPublicApi()
}

func testEmail() {
	email.Start()
	email.SendEmail(&machinepb.Email{
		Subject:   "maintain me",
		Body:      "somehting broked",
		Recipient: machinepb.EmailRecipient_EMAIL_RECIPIENT_MAINTENANCE,
	})
}

func printProfiles() {
	for i, v := range vialprofiles.GetSystemVialConfigurationSnapshot().Profiles {
		fmt.Println(i)
		fmt.Println(v)
	}
}

// serves the public api locally, with control enabled
func testPublicApi() {
	mqtt.Start(config.BrokerHost())
	sm := session.NewSessionManager(false)

	api := publicapi.New()
	events.Start(sm, api)
	vialprofiles.Start(sm, api)
	control.Start(api)
	control.SetEnabled(true)

	if err := api.Start(publicapi.Options{Addr: "127.0.0.1:8789"}); err != nil {
		panic(err)
	}
	select {}
}
