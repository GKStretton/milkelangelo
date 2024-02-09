package main

import (
	"fmt"
	"time"

	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/dark/services/goo/actor"
	"github.com/gkstretton/dark/services/goo/actor/decider"
	"github.com/gkstretton/dark/services/goo/ebsinterface"
	"github.com/gkstretton/dark/services/goo/email"
	"github.com/gkstretton/dark/services/goo/events"
	"github.com/gkstretton/dark/services/goo/mqtt"
	"github.com/gkstretton/dark/services/goo/session"
	"github.com/gkstretton/dark/services/goo/twitchapi"
	"github.com/gkstretton/dark/services/goo/types"
	"github.com/gkstretton/dark/services/goo/util"
	"github.com/gkstretton/dark/services/goo/vialprofiles"
)

// tests for human verification during development
func runAdHocTests() {
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

func testRandomVialPos() {
	i := 0
	for {
		if i > 100 {
			break
		}
		i++

		fmt.Print(decider.GetRandomVialPos())
		fmt.Print(", ")
	}
}

func testSampleUnitCircle() {
	fmt.Println(util.SampleRandomUnitCircleCoordinate())
}

func testActor() {
	mqtt.Start()
	sm := session.NewSessionManager(false)
	events.Start(sm)
	twitchApi := twitchapi.Start()

	actor.LaunchActor(twitchApi, 3*time.Minute)
}

// subscribes to ebs and twitch chat votes and prints the received votes
func testEBSAndChatVoting() {
	mqtt.Start()
	sm := session.NewSessionManager(false)
	events.Start(sm)

	ebs, err := ebsinterface.NewExtensionSession(time.Hour * 2)
	if err != nil {
		panic(err)
	}
	ebs.UpdateCurrentVoteStatus(&types.VoteStatus{
		VoteType: types.VoteTypeCollection,
		CollectionVoteStatus: &types.CollectionVoteStatus{
			TotalVotes: 5,
			VoteCounts: map[uint64]uint64{5: 25},
		},
	})
	ebsCh := ebs.SubscribeVotes()
	defer ebs.UnsubscribeVotes(ebsCh)
	api := twitchapi.Start()
	chatCh, unsub := api.SubscribeChatVotes(types.VoteTypeLocation)
	defer unsub()
	chatCh2, unsub2 := api.SubscribeChatVotes(types.VoteTypeCollection)
	defer unsub2()

	for {
		select {
		case vote := <-ebsCh:
			fmt.Printf("got ebs vote:\n%+v\n\n", vote)
		case vote := <-chatCh:
			fmt.Printf("got location chat vote:\n%+v\n\n", vote)
		case vote := <-chatCh2:
			fmt.Printf("got collection chat vote:\n%+v\n\n", vote)
		}
	}
}
