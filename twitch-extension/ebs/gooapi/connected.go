package gooapi

import (
	"fmt"
	"net/http"
	"sync"
)

type connectedGooApi struct {
	internalSecret []byte
	listenAddr     string

	subs     []chan *message
	subsLock sync.Mutex

	stateUpdateCallback func(state GooStateUpdate)
}

func NewConnectedGooApi(sharedSecret string, listenAddr string) (*connectedGooApi, error) {
	if sharedSecret == "" {
		return nil, fmt.Errorf("shared secret not set for goo")
	}
	g := &connectedGooApi{
		internalSecret: []byte(sharedSecret),
		listenAddr:     listenAddr,
		subs:           []chan *message{},
		subsLock:       sync.Mutex{},
	}

	http.HandleFunc("/listen", g.listenHandler)
	http.HandleFunc("/update-state", g.updateHandler)

	return g, nil
}

func (g *connectedGooApi) Start() {
	err := http.ListenAndServe(g.listenAddr, nil)
	if err != nil {
		l.Fatalf("error in listen and server for goo requests", err)
	}
}

func (g *connectedGooApi) SetStateUpdateCallback(callback func(state GooStateUpdate)) {
	g.stateUpdateCallback = callback
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

func (g *connectedGooApi) Dispense(x, y float32) error {
	m := dispenseRequest{
		X: x,
		Y: y,
	}

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

func (g *connectedGooApi) ReportEbsState(state EbsStateReport) error {
	return g.sendMessage(&message{
		MessageType: stateReportType,
		Data:        state,
	})
}
