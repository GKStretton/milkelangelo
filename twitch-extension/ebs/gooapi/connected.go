package gooapi

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gkstretton/study-of-light/twitch-ebs/common"
)

type connectedGooApi struct {
	internalSecret []byte
	listenAddr     string

	subs     []chan *message
	subsLock sync.Mutex
}

func NewConnectedGooApi(internalSecretPath string, listenAddr string) (*connectedGooApi, error) {
	internalSecret, err := common.GetSecret(internalSecretPath)
	if err != nil {
		return nil, fmt.Errorf("could not load internal shared secret: %w", err)
	}

	g := &connectedGooApi{
		internalSecret: internalSecret,
		listenAddr:     listenAddr,
		subs:           []chan *message{},
		subsLock:       sync.Mutex{},
	}

	http.HandleFunc("/listen", g.listenHandler)

	return g, nil
}

func (g *connectedGooApi) Start() {
	err := http.ListenAndServe(g.listenAddr, nil)
	if err != nil {
		l.Fatalf("error in listen and server for goo requests", err)
	}
}

func (g *connectedGooApi) CollectFromVial(vial int) error {
	m := collectionRequest{
		Id: vial,
	}

	return g.sendMessage(&message{
		MessageType: collectionRequestType,
		Data:        m,
	})
}

func (g *connectedGooApi) Dispense() error {
	m := dispenseRequest{}

	return g.sendMessage(&message{
		MessageType: dispenseRequestType,
		Data:        m,
	})
}

func (g *connectedGooApi) GoToPosition(x, y float32) error {
	m := goToRequest{
		X: x,
		Y: y,
	}

	return g.sendMessage(&message{
		MessageType: goToRequestType,
		Data:        m,
	})
}
