package ebsinterface

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gkstretton/dark/services/goo/types"
	"github.com/r3labs/sse"
)

const ebsListenUrl = "http://localhost:8080/listen"

var subs = []chan *types.Vote{}
var subsLock = sync.Mutex{}

func (e *ExtensionSession) SubscribeVotes() <-chan *types.Vote {
	subsLock.Lock()
	defer subsLock.Unlock()

	c := make(chan *types.Vote)
	subs = append(subs, c)
	return c
}

func (e *ExtensionSession) UnsubscribeVotes(c <-chan *types.Vote) {
	subsLock.Lock()
	defer subsLock.Unlock()

	for i, sub := range subs {
		if sub == c {
			subs = append(subs[:i], subs[i+1:]...)
			return
		}
	}
}

func (e *ExtensionSession) distributeVote(data *types.Vote) {
	subsLock.Lock()
	defer subsLock.Unlock()

	for _, sub := range subs {
		select {
		case sub <- data:
		default:
		}
	}
}

// rawVote is received over sse channel
type rawVote struct {
	Data          []byte
	OpaqueUserID  string
	IsBroadcaster bool
}

// connect listens to the ebs vote stream
func (e *ExtensionSession) connect() error {
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
			vote := &rawVote{}
			err := json.Unmarshal(event.Data, vote)
			if err != nil {
				l.Printf("failed to unmarshal rawVote: %v\n", err)
				continue
			}

			voteDetails := &types.VoteDetails{}
			err = json.Unmarshal(vote.Data, voteDetails)
			if err != nil {
				l.Printf("failed to unmarshal voteDetails: %v\n", err)
				continue
			}

			if err := validateVote(voteDetails); err != nil {
				l.Printf("invalid vote: %v\n", err)
				continue
			}

			e.distributeVote(&types.Vote{
				Data:          *voteDetails,
				OpaqueUserID:  vote.OpaqueUserID,
				IsBroadcaster: vote.IsBroadcaster,
			})
		}
	}()

	go func() {
		<-e.exitCh
		client.Unsubscribe(ch)
		close(ch)
	}()

	return nil
}

func validateVote(d *types.VoteDetails) error {
	if d == nil {
		return fmt.Errorf("nil vote")
	}
	if d.VoteType == types.VoteTypeCollection && d.CollectionVote == nil {
		return fmt.Errorf("nil collection vote")
	}
	if d.VoteType == types.VoteTypeLocation && d.LocationVote == nil {
		return fmt.Errorf("nil location vote")
	}

	if d.VoteType == types.VoteTypeLocation {
		v := d.LocationVote
		if v.N > 11 {
			return fmt.Errorf("brain too big")
		}
		if v.X < -1 || v.X > 1 {
			return fmt.Errorf("x out of range -1 - 1")
		}
		if v.Y < -1 || v.Y > 1 {
			return fmt.Errorf("y out of range -1 - 1")
		}
	}
	return nil
}
