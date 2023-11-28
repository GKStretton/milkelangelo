package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
)

var viewerRequestLogging = flag.Bool("requestLogging", false, "if enabled, log viewer requests")

func voteHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	claims, err := verifyUserRequest(r)
	if err != nil {
		fmt.Printf("error verifying user request: %v\n", err)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	data, err := io.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		fmt.Printf("failed to read body: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if *viewerRequestLogging {
		fmt.Printf("got vote data from %s\n", r.RemoteAddr)
	}
	err = sendVote(&vote{
		data:          data,
		opaqueUserID:  claims.OpaqueUserID,
		isBroadcaster: claims.Role != "broadcaster",
	})
	if err != nil {
		fmt.Fprintf(w, "failed to vote: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
