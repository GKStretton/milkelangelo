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

type extensionSession struct {
	broadcastToken    string
	ebsListeningToken string
	cleanUpDone       bool
	lock              sync.Mutex

	// used to disconnect from ebs
	exitCh chan struct{}
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
		exitCh:            make(chan struct{}),
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
