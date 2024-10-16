package main

import (
	"fmt"
	"time"

	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/dark/services/goo/actor"
	"github.com/gkstretton/dark/services/goo/ebsinterface"
	"github.com/gkstretton/dark/services/goo/email"
	"github.com/gkstretton/dark/services/goo/events"
	"github.com/gkstretton/dark/services/goo/mqtt"
	"github.com/gkstretton/dark/services/goo/session"
	"github.com/gkstretton/dark/services/goo/twitchapi"
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

func testActor() {
	mqtt.Start()
	sm := session.NewSessionManager(false)
	events.Start(sm)
	twitchApi := twitchapi.Start()

	actor.LaunchActor(twitchApi, 3*time.Minute, 1, true)
}

// subscribes to ebs and twitch chat votes and prints the received votes
func testEBS() {
	mqtt.Start()
	sm := session.NewSessionManager(false)
	events.Start(sm)

	ebs, err := ebsinterface.NewExtensionSession("http://localhost:8788", time.Hour*2)
	if err != nil {
		panic(err)
	}
	defer ebs.CleanUp()

	ebsCh := ebs.SubscribeVotes()
	defer ebs.UnsubscribeVotes(ebsCh)

	for message := range ebsCh {
		fmt.Printf("got ebs message '%s':\n\t%+v\n\t%+v\n\t%+v\n\n", message.Type, message.DispenseRequest, message.CollectionRequest, message.GoToRequest)
	}
}
