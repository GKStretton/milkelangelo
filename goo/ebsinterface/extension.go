package ebsinterface

import (
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gkstretton/dark/services/goo/types"
)

var l = log.New(os.Stdout, "[EBS Interface] ", log.Flags())

type extensionSession struct {
	ebsToken string

	// used to disconnect from ebs
	exitCh chan struct{}

	subs     []chan *types.EbsMessage
	subsLock sync.Mutex

	ebsAddress string

	gooStateLock sync.Mutex
	gooState     types.GooState
	gooStateChan chan struct{}

	ebsStateLock sync.Mutex
	ebsState     *types.EbsStateReport
}

type EbsApi interface {
	SubscribeMessages() <-chan *types.EbsMessage
	UnsubscribeMessages(c <-chan *types.EbsMessage)
	UpdateState(func(state *types.GooState))
	GetEbsState() *types.EbsStateReport
}

func NewExtensionSession(ebsAddress string) (*extensionSession, error) {
	http.DefaultClient.Timeout = 5 * time.Second

	elt, err := getEBSListeningToken()
	if err != nil {
		return nil, err
	}

	es := &extensionSession{
		ebsToken:     elt,
		exitCh:       make(chan struct{}),
		subs:         []chan *types.EbsMessage{},
		ebsAddress:   ebsAddress,
		gooStateChan: make(chan struct{}, 1),
	}

	l.Println("connecting to ebs...")
	err = es.connect()
	if err != nil {
		return nil, err
	}

	// send goo state update to ebs whenever there's a change
	go es.stateSender()

	return es, nil
}
