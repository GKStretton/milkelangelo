package ebsinterface

import (
	"errors"
	"flag"
	"log"
	"os"
	"sync"
	"time"
)

var l = log.New(os.Stdout, "[EBS Interface] ", log.Flags())

const channelId = "807784320"

var (
	useLocalEBS = flag.Bool("useLocalEBS", true, "if true, use local ebs rather than hosted one")
)

type ExtensionSession struct {
	broadcastToken    string
	ebsListeningToken string
	cleanUpDone       bool
	lock              sync.Mutex

	// used to disconnect from ebs
	exitCh             chan struct{}
	triggerBroadcast   chan struct{}
	broadcastDataCache *broadcastData
}

func NewExtensionSession(dur time.Duration) (*ExtensionSession, error) {
	bt, err := getBroadcastToken(dur)
	if err != nil {
		return nil, err
	}
	elt, err := getEBSListeningToken(dur)
	if err != nil {
		return nil, err
	}

	es := &ExtensionSession{
		broadcastToken:     bt,
		ebsListeningToken:  elt,
		exitCh:             make(chan struct{}),
		triggerBroadcast:   make(chan struct{}),
		broadcastDataCache: &broadcastData{},
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

	go es.regularRobotStatusUpdate()
	go es.regularBroadcast()

	return es, nil
}

// trigger running of the EBS on fly.io
func (e *ExtensionSession) launch() error {
	if *useLocalEBS {
		return nil
	}
	// todo: implement
	return errors.New("launch not implemented for hosted ebs")
}

// CleanUp will be automatically called after duration specified on creation.
// If exiting early, call manually.
func (e *ExtensionSession) CleanUp() {
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
