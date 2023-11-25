package twitchextension

import (
	"errors"
	"flag"
	"fmt"
	"sync"
	"time"
)

const channelId = "807784320"

var (
	useLocalEBS = flag.Bool("useLocalEBS", true, "if true, use local ebs rather than hosted one")
)

type extensionSession struct {
	broadcastToken    string
	ebsListeningToken string
	cleanUpDone       bool
	lock              sync.Mutex
}

func NewExtensionSession(dur time.Duration) (*extensionSession, error) {
	bt, err := getBroadcastToken(dur)
	if err != nil {
		return nil, err
	}
	elt, err := getEBSListeningToken(dur)
	if err != nil {
		return nil, err
	}

	es := &extensionSession{
		broadcastToken:    bt,
		ebsListeningToken: elt,
	}

	err = es.launch()
	if err != nil {
		return nil, err
	}
	go func() {
		time.Sleep(dur)
		es.CleanUp()
	}()
	return es, nil
}

// trigger running of the EBS on fly.io
func (e *extensionSession) launch() error {
	if *useLocalEBS {
		return nil
	}
	// todo: implement
	return errors.New("launch not implemented for hosted ebs")
}

func (e *extensionSession) CleanUp() {
	e.lock.Lock()
	defer e.lock.Unlock()
	if e.cleanUpDone {
		return
	}
	defer func() { e.cleanUpDone = true }()

	if *useLocalEBS {
		return
	}
	// todo: implement
	fmt.Println("error, cleanup not implemented for hosted ebs")
}
