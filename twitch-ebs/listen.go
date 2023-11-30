package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/golang-jwt/jwt"
)

type Vote struct {
	Data          []byte
	OpaqueUserID  string
	IsBroadcaster bool
}

var subs = []chan *Vote{}
var subsLock = sync.Mutex{}

func registerVoteListener(c chan *Vote) {
	subsLock.Lock()
	defer subsLock.Unlock()

	subs = append(subs, c)
}

func removeVoteListener(c chan *Vote) {
	subsLock.Lock()
	defer subsLock.Unlock()

	for i, sub := range subs {
		if sub == c {
			subs = append(subs[:i], subs[i+1:]...)
			close(c)
			break
		}
	}
}

func sendVote(v *Vote) error {
	subsLock.Lock()
	defer subsLock.Unlock()

	if len(subs) == 0 {
		return errors.New("nobody's listening...")
	}

	for _, sub := range subs {
		select {
		case sub <- v:
		default:
		}
	}
	return nil
}

func listenHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	claims, err := verifyInternalRequest(r)
	if err != nil {
		httpErr(&w, http.StatusUnauthorized, "failed to verify internal (goo) request: %v", err)
		return
	}
	PrintRequesterInfo(r, claims)

	w.WriteHeader(http.StatusOK)
	w.(http.Flusher).Flush()

	// 1 buffer
	c := make(chan *Vote, 1)
	// c subscribes to vote stream from viewers
	registerVoteListener(c)
	defer removeVoteListener(c)

	// do the listening and returning
	for {
		select {
		case <-r.Context().Done():
			return // exit handler, remove listener
		case vote := <-c:
			data, err := json.Marshal(vote)
			if err != nil {
				fmt.Printf("failed to marshal vote: %v\n", err)
				continue
			}
			fmt.Fprintf(w, "event: vote\n")
			fmt.Fprintf(w, "data: %s\n\n", data)
			w.(http.Flusher).Flush()
		}
	}
}

// PrintRequesterInfo prints information about the requester from an http.Request
func PrintRequesterInfo(r *http.Request, claims *jwt.StandardClaims) {
	fmt.Printf("verified new listener request from %s\n", r.RemoteAddr)

	fmt.Println("Headers:")
	for name, headers := range r.Header {
		for _, h := range headers {
			fmt.Printf("\t%v: %v\n", name, h)
		}
	}
	fmt.Printf("claims: %+v\n\n\n", claims)
}
