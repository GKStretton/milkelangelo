package gooapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
)

type messageType string

const (
	dispenseRequestType   messageType = "dispense"
	collectionRequestType messageType = "collection"
	goToRequestType       messageType = "goto"
)

type dispenseRequest struct {
	X float32
	Y float32
}

type collectionRequest struct {
	Id int
}

type goToRequest struct {
	X float32
	Y float32
}

type message struct {
	MessageType messageType
	Data        interface{}
}

func (g *connectedGooApi) registerMessageListener(c chan *message) {
	g.subsLock.Lock()
	defer g.subsLock.Unlock()

	g.subs = append(g.subs, c)

	l.Debug("registered new message listener (i.e. goo instance)")
}

func (g *connectedGooApi) removeMessageListener(c chan *message) {
	g.subsLock.Lock()
	defer g.subsLock.Unlock()

	for i, sub := range g.subs {
		if sub == c {
			g.subs = append(g.subs[:i], g.subs[i+1:]...)
			close(c)
			break
		}
	}

	l.Debug("removed a message listener (i.e. goo instance)")
}

func (g *connectedGooApi) sendMessage(msg *message) error {
	g.subsLock.Lock()
	defer g.subsLock.Unlock()

	if len(g.subs) == 0 {
		return errors.New("no listener (i.e. no goo instance)")
	}

	for _, sub := range g.subs {
		select {
		case sub <- msg:
		default:
		}
	}

	l.Debugf("sent message to %d connected goo instance(s): %s - %+v", len(g.subs), msg.MessageType, msg.Data)

	return nil
}

func (g *connectedGooApi) listenHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "chunked")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	claims, err := g.verifyInternalRequest(r)
	if err != nil {
		httpErr(&w, http.StatusUnauthorized, "failed to verify internal (goo) request: %v", err)
		return
	}
	PrintRequesterInfo(r, claims)

	w.WriteHeader(http.StatusOK)
	w.(http.Flusher).Flush()

	// 1 buffer
	c := make(chan *message, 1)
	// c subscribes to message stream from any twitch clients
	g.registerMessageListener(c)
	defer g.removeMessageListener(c)

	// do the listening and returning
	for {
		select {
		case <-r.Context().Done():
			return // exit handler, remove listener
		case msg := <-c:
			data, err := json.Marshal(msg.Data)
			if err != nil {
				fmt.Printf("failed to marshal %s message: %v\n", msg.MessageType, err)
				continue
			}
			fmt.Fprintf(w, "event: %s\n", msg.MessageType)
			fmt.Fprintf(w, "data: %s\n\n", data)
			w.(http.Flusher).Flush()
		}
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET")
	(*w).Header().Set("Access-Control-Allow-Headers", "Authorization")
}

// PrintRequesterInfo prints information about the requester from an http.Request
func PrintRequesterInfo(r *http.Request, claims *jwt.StandardClaims) {
	if claims == nil {
		return
	}
	l.Infof("verified new listener request from %s\n", r.RemoteAddr)

	l.Infof("Headers:")
	for name, headers := range r.Header {
		for _, h := range headers {
			l.Infof("\t%v: %v", name, h)
		}
	}
	l.Infof("claims: %+v\n\n", claims)
}

func httpErr(w *http.ResponseWriter, code int, s string, args ...interface{}) {
	msg := fmt.Errorf(s, args...)
	(*w).WriteHeader(code)
	fmt.Fprintln(*w, msg)
	l.Errorf("%d: %s", code, msg)
}
