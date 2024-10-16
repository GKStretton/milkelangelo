package twitchapi

import (
	"github.com/gempir/go-twitch-irc/v4"
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
