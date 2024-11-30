package ebsinterface

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/gkstretton/dark/services/goo/types"
	"github.com/r3labs/sse"
)

func (e *extensionSession) SubscribeMessages() <-chan *types.EbsMessage {
	e.subsLock.Lock()
	defer e.subsLock.Unlock()

	c := make(chan *types.EbsMessage)
	e.subs = append(e.subs, c)
	return c
}

func (e *extensionSession) UnsubscribeMessages(c <-chan *types.EbsMessage) {
	e.subsLock.Lock()
	defer e.subsLock.Unlock()

	for i, sub := range e.subs {
		if sub == c {
			e.subs = append(e.subs[:i], e.subs[i+1:]...)
			return
		}
	}
}

func (e *extensionSession) distributeMessage(msg *types.EbsMessage) {
	e.subsLock.Lock()
	defer e.subsLock.Unlock()

	for _, sub := range e.subs {
		select {
		case sub <- msg:
		default:
		}
	}
}

// connect listens to the ebs message stream
func (e *extensionSession) connect() error {
	result, err := url.JoinPath(e.ebsAddress, "/listen")
	if err != nil {
		return fmt.Errorf("error forming ebs listen url: %s", err)
	}
	client := sse.NewClient(result)
	client.Headers["Authorization"] = "Bearer " + e.ebsListeningToken
	// client.ReconnectStrategy = &backoff.StopBackOff{}
	client.ReconnectNotify = func(err error, d time.Duration) {
		l.Println(err)
	}

	ch := make(chan *sse.Event)
	client.OnDisconnect(func(c *sse.Client) {
		l.Printf("disconnected: %v\n", c)
	})
	err = client.SubscribeChan("", ch)
	if err != nil {
		return fmt.Errorf("failed to subscribe to twitch ebs sse: %v", err)
	}
	go func() {
		for {
			event, ok := <-ch
			if !ok {
				l.Println("closing twitch ebs message listener")
				return
			}

			msg := &types.EbsMessage{
				Type: types.EbsMessageType(event.Event),
			}

			switch msg.Type {
			case types.EbsDispenseRequest:
				l.Printf("got dispense request from ebs")
				err := json.Unmarshal(event.Data, &msg.DispenseRequest)
				if err != nil {
					l.Printf("error unmarshalling dispense request from ebs: %s", err)
					continue
				}
			case types.EbsCollectionRequest:
				l.Printf("got collection request from ebs")
				err := json.Unmarshal(event.Data, &msg.CollectionRequest)
				if err != nil {
					l.Printf("error unmarshalling collection request from ebs: %s", err)
					continue
				}
			case types.EbsGoToRequest:
				l.Printf("got goto request from ebs")
				err := json.Unmarshal(event.Data, &msg.GoToRequest)
				if err != nil {
					l.Printf("error unmarshalling goto request from ebs: %s", err)
					continue
				}
			default:
				l.Printf("unrecognised ebs message type '%s'", msg.Type)
				continue
			}

			e.distributeMessage(msg)
		}
	}()

	go func() {
		<-e.exitCh
		client.Unsubscribe(ch)
		close(ch)
	}()

	return nil
}
