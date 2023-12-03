package main

import (
	"fmt"
	"time"

	"github.com/gkstretton/dark/services/goo/ebsinterface"
	"github.com/gkstretton/dark/services/goo/twitchapi"
	"github.com/gkstretton/dark/services/goo/types"
)

// tests for human verification during development
func runAdHocTests() {
	testActor()
}

func testActor() {
	fmt.Println("todo")
}

// subscribes to ebs and twitch chat votes and prints the received votes
func testEBSAndChatVoting() {
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
