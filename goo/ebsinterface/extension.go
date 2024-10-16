package ebsinterface

import (
	"errors"
	"flag"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gkstretton/dark/services/goo/types"
)

var l = log.New(os.Stdout, "[EBS Interface] ", log.Flags())

var (
	useLocalEBS = flag.Bool("useLocalEBS", true, "if true, use local ebs rather than hosted one")
)

type extensionSession struct {
	ebsListeningToken string
	cleanUpDone       bool
	lock              sync.Mutex

	// used to disconnect from ebs
	exitCh           chan struct{}
	triggerBroadcast chan struct{}

	subs     []chan *types.EbsMessage
	subsLock sync.Mutex

	ebsAddress string
}

type EbsApi interface {
	SubscribeMessages() <-chan *types.EbsMessage
	UnsubscribeMessages(c <-chan *types.EbsMessage)
}

func NewExtensionSession(ebsAddress string, dur time.Duration) (*extensionSession, error) {
	elt, err := getEBSListeningToken(dur)
	if err != nil {
		return nil, err
	}

	es := &extensionSession{
		ebsListeningToken: elt,
		exitCh:            make(chan struct{}),
		triggerBroadcast:  make(chan struct{}),
		subs:              []chan *types.EbsMessage{},
		ebsAddress:        ebsAddress,
	}

	err = es.launch()
	if err != nil {
		return nil, err
	}
	go func() {
		time.Sleep(dur)
		es.CleanUp()
	}()

	l.Println("connecting to ebs...")
	err = es.connect()
	if err != nil {
		es.CleanUp()
		return nil, err
	}
	l.Println("connected to ebs.")

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

// CleanUp will be automatically called after duration specified on creation.
// If exiting early, call manually.
func (e *extensionSession) CleanUp() {
	e.lock.Lock()
	defer e.lock.Unlock()
	if e.cleanUpDone {
		return
	}
	defer func() { e.cleanUpDone = true }()

	close(e.exitCh)

	if *useLocalEBS {
		return
	}
	// todo: implement
	l.Println("error, cleanup not implemented for hosted ebs")
}
