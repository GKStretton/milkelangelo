// Package publicapi is the HTTP API remote clients use to control the machine.
//
// It does no authentication: who may reach it is decided by whatever is in
// front of it (e.g. a reverse proxy or tunnel with its own login). It decides
// who has control right now (a claim lease) and serves state to clients. What
// the machine may do is decided by the control package.
package publicapi

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gkstretton/dark/services/goo/types"
)

var l = log.New(os.Stdout, "[publicapi] ", log.Flags())

// Api holds the state served to clients, and who currently has control.
type Api struct {
	lock     sync.Mutex
	gooState types.GooState

	claimToken  string
	expiryTimer *time.Timer

	subs map[chan struct{}]struct{}
}

// state is sent to clients by /state and /events
type state struct {
	GooState *types.GooState
	// Claimed is true while a client holds control
	Claimed bool
}

func New() *Api {
	return &Api{
		subs: map[chan struct{}]struct{}{},
	}
}

// UpdateState implements types.GooStateUpdater
func (a *Api) UpdateState(f func(state *types.GooState)) {
	a.lock.Lock()
	defer a.lock.Unlock()

	f(&a.gooState)
	a.notify()
}

func (a *Api) marshalState() ([]byte, error) {
	a.lock.Lock()
	defer a.lock.Unlock()

	return json.Marshal(state{
		GooState: &a.gooState,
		Claimed:  a.claimToken != "",
	})
}

// subscribe returns a channel that receives a value whenever state changes.
// Changes are coalesced, so a slow subscriber only sees the latest state.
func (a *Api) subscribe() chan struct{} {
	a.lock.Lock()
	defer a.lock.Unlock()

	c := make(chan struct{}, 1)
	a.subs[c] = struct{}{}
	return c
}

func (a *Api) unsubscribe(c chan struct{}) {
	a.lock.Lock()
	defer a.lock.Unlock()

	delete(a.subs, c)
}

// notify must be called with a.lock held
func (a *Api) notify() {
	for c := range a.subs {
		select {
		case c <- struct{}{}:
		default:
		}
	}
}
