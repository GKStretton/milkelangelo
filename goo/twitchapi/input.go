package twitchapi

import (
	"fmt"
	"strings"

	"github.com/gempir/go-twitch-irc/v4"
	"github.com/gkstretton/dark/services/goo/types"
	"github.com/gkstretton/dark/services/goo/vialprofiles"
)

type Message struct {
	twitch.PrivateMessage
}

func (m *Message) IsSelf() bool {
	return m.User.ID == channelId
}

func (api *TwitchApi) SubscribeChat() chan *Message {
	api.lock.Lock()
	defer api.lock.Unlock()
	c := make(chan *Message, 10)
	api.subs = append(api.subs, c)
	return c
}

func (api *TwitchApi) UnsubscribeChat(c chan *Message) {
	api.lock.Lock()
	defer api.lock.Unlock()
	for i, sub := range api.subs {
		if sub == c {
			api.subs = append(api.subs[:i], api.subs[i+1:]...)
			close(c)
			break
		}
	}
}

func (api *TwitchApi) setupHandlers() {
	api.ircClient.OnPrivateMessage(func(message twitch.PrivateMessage) {
		api.lock.Lock()
		// send message to all subscribers
		for _, c := range api.subs {
			select {
			case c <- &Message{message}:
			default:
			}
		}
		api.lock.Unlock()
	})
}

func (api *TwitchApi) SubscribeChatVotes(voteType types.VoteType) (voteCh chan *types.Vote, unsubscribe func()) {
	_, vialPosToName := vialprofiles.GetVialOptionsAndMap()

	msgCh := api.SubscribeChat()

	voteCh = make(chan *types.Vote)
	go func() {
		for {
			msg, ok := <-msgCh
			if !ok {
				close(voteCh)
				return
			}
			vote, err := TwitchMessageToVote(voteType, msg, vialPosToName)
			if err != nil {
				fmt.Printf("failed to parse vote from %s: %s\n", msg.Message, err)
				api.Reply(msg.ID, err.Error())
				continue
			}
			if vote == nil {
				continue
			}
			if vote.Data.LocationVote != nil && vote.Data.LocationVote.N > 2 {
				n := vote.Data.LocationVote.N
				api.Reply(msg.ID, fmt.Sprintf("%dD%s", n, strings.Repeat("!?", int(n-2))))
			}
			voteCh <- vote
		}
	}()

	return voteCh, func() {
		api.UnsubscribeChat(msgCh)
	}
}
