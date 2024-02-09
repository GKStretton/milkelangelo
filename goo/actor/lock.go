package actor

import (
	"fmt"
	"sync"

	"github.com/gkstretton/asol-protos/go/topics_backend"
	"github.com/gkstretton/dark/services/goo/mqtt"
)

type ActorLock struct {
	isRunning bool
	lock      sync.Mutex
}

func (a *ActorLock) Set(isRunning bool) {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.isRunning = isRunning
	mqtt.Publish(topics_backend.TOPIC_ACTOR_STATUS_RESP, fmt.Sprintf("%t", a.isRunning))
}

func (a *ActorLock) Get() bool {
	a.lock.Lock()
	defer a.lock.Unlock()

	return a.isRunning
}
