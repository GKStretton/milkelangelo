package ebsinterface

import (
	"fmt"
	"sync"
	"time"

	"github.com/r3labs/sse"
)

const ebsListenUrl = "http://localhost:8080/listen"

// json string
type VoteData = string

var subs = []chan VoteData{}
var subsLock = sync.Mutex{}

func (e *extensionSession) SubscribeVotes() <-chan VoteData {
	subsLock.Lock()
	defer subsLock.Unlock()

	c := make(chan VoteData)
	subs = append(subs, c)
	return c
}

func (e *extensionSession) UnsubscribeVotes(c <-chan VoteData) {
	subsLock.Lock()
	defer subsLock.Unlock()

	for i, sub := range subs {
		if sub == c {
			subs = append(subs[:i], subs[i+1:]...)
			return
		}
	}
}

func (e *extensionSession) distributeVote(data VoteData) {
	l.Printf("got vote: %s\n", data)
	subsLock.Lock()
	defer subsLock.Unlock()

	for _, sub := range subs {
		select {
		case sub <- data:
		default:
		}
	}
}

// connect listens to the ebs vote stream
func (e *extensionSession) connect() error {
	client := sse.NewClient(ebsListenUrl)
	client.Headers["Authorization"] = "Bearer " + e.ebsListeningToken
	// client.ReconnectStrategy = &backoff.StopBackOff{}
	client.ReconnectNotify = func(err error, d time.Duration) {
		l.Println(err)
	}

	ch := make(chan *sse.Event)
	client.OnDisconnect(func(c *sse.Client) {
		l.Printf("disconnected: %v\n", c)
	})
	err := client.SubscribeChan("", ch)
	if err != nil {
		return fmt.Errorf("failed to subscribe to twitch ebs sse: %v", err)
	}
	go func() {
		for {
			event, ok := <-ch
			if !ok {
				l.Println("closing twitch ebs vote listener")
				return
			}
			//todo: rate limit to prevent ddos, here or in ebs?
			e.distributeVote(VoteData(event.Data))
		}
	}()

	go func() {
		<-e.exitCh
		close(ch)
	}()

	return nil
}
