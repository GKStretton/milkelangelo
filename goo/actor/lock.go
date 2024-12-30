package actor

import (
	"fmt"
	"sync"

	"github.com/gkstretton/asol-protos/go/topics_backend"
	"github.com/gkstretton/dark/services/goo/ebsinterface"
	"github.com/gkstretton/dark/services/goo/mqtt"
	"github.com/gkstretton/dark/services/goo/types"
)

type ActorLock struct {
	isRunning bool
	lock      sync.Mutex
}

func (a *ActorLock) Set(ebsApi ebsinterface.EbsApi, isRunning bool) {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.isRunning = isRunning
	mqtt.Publish(topics_backend.TOPIC_ACTOR_STATUS_RESP, fmt.Sprintf("%t", a.isRunning))

	if ebsApi != nil {
		ebsApi.UpdateState(func(state *types.GooState) {
			state.ActorRunning = isRunning
		})
	}
}

func (a *ActorLock) Get() bool {
	a.lock.Lock()
	defer a.lock.Unlock()

	return a.isRunning
}
