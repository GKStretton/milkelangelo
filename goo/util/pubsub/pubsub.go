package pubsub

import (
	"errors"
	"sync"
)

type EventDistributor[T any] struct {
	lock            sync.Mutex
	subs            []chan T
	senderConnected bool
	senderChan      chan T
}

func (ed *EventDistributor[T]) RegisterListener() (<-chan T, error) {
	ed.lock.Lock()
	defer ed.lock.Unlock()

	if !ed.senderConnected {
		return nil, errors.New("sender not connected")
	}

	c := make(chan T)
	ed.subs = append(ed.subs, c)
	return c, nil
}

func (ed *EventDistributor[T]) DeregisterListener(c <-chan T) {
	ed.lock.Lock()
	defer ed.lock.Unlock()

	for i, sub := range ed.subs {
		if sub == c {
			ed.subs = append(ed.subs[:i], ed.subs[i+1:]...)
		}
	}
}

// todo: write sender stuff
