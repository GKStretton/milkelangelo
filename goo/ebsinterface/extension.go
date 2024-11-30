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
	ebsListeningToken string

	// used to disconnect from ebs
	exitCh chan struct{}

	subs     []chan *types.EbsMessage
	subsLock sync.Mutex

	ebsAddress string

	gooStateLock sync.Mutex
	gooState     types.GooState
}

type EbsApi interface {
	SubscribeMessages() <-chan *types.EbsMessage
	UnsubscribeMessages(c <-chan *types.EbsMessage)
	UpdateState(func(state *types.GooState))
}

func NewEbsApi(ebsAddress string) *extensionSession {
	panic("NewEbsApi not implemented, only NewExtensionSession(duration) currently")
}

func NewExtensionSession(ebsAddress string) (*extensionSession, error) {
	http.DefaultClient.Timeout = 5 * time.Second

	elt, err := getEBSListeningToken()
	if err != nil {
		return nil, err
	}

	es := &extensionSession{
		ebsListeningToken: elt,
		exitCh:            make(chan struct{}),
		subs:              []chan *types.EbsMessage{},
		ebsAddress:        ebsAddress,
	}

	l.Println("connecting to ebs...")
	err = es.connect()
	if err != nil {
		return nil, err
	}
	l.Println("connected to ebs.")

	return es, nil
}
