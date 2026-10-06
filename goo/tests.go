package main

import (
	"fmt"

	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/dark/services/goo/config"
	"github.com/gkstretton/dark/services/goo/control"
	"github.com/gkstretton/dark/services/goo/ebsinterface"
	"github.com/gkstretton/dark/services/goo/email"
	"github.com/gkstretton/dark/services/goo/events"
	"github.com/gkstretton/dark/services/goo/mqtt"
	"github.com/gkstretton/dark/services/goo/session"
	"github.com/gkstretton/dark/services/goo/vialprofiles"
)

// tests for human verification during development
func runAdHocTests() {
	testEBS()
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

// connects to a local ebs and enables control, printing received messages
func testEBS() {
	mqtt.Start(config.BrokerHost())
	sm := session.NewSessionManager(false)

	ebs, err := ebsinterface.NewExtensionSession("http://localhost:8788")
	if err != nil {
		panic(err)
	}
	events.Start(sm, ebs)
	vialprofiles.Start(sm, ebs)
	control.Start(ebs)
	control.SetEnabled(true)

	ebsCh := ebs.SubscribeMessages()
	defer ebs.UnsubscribeMessages(ebsCh)

	for message := range ebsCh {
		fmt.Printf("got ebs message '%s':\n\t%+v\n\t%+v\n\t%+v\n\n", message.Type, message.DispenseRequest, message.CollectionRequest, message.GoToRequest)
	}
}
