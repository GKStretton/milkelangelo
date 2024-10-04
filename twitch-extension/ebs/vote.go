package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
)

var viewerRequestLogging = flag.Bool("requestLogging", true, "if enabled, log viewer requests")

func voteHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	claims, err := verifyUserRequest(r)
	if err != nil {
		httpErr(&w, http.StatusForbidden, "error verifying user request: %v", err)
		return
	}

	data, err := io.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		httpErr(&w, http.StatusBadRequest, "failed to read body: %v", err)
		return
	}
	if *viewerRequestLogging {
		fmt.Printf("got vote data from %s: %s\n", r.RemoteAddr, string(data))
	}
	err = sendVote(&Vote{
		Data:          data,
		OpaqueUserID:  claims.OpaqueUserID,
		IsBroadcaster: claims.Role != "broadcaster",
	})
	if err != nil {
		httpErr(&w, http.StatusInternalServerError, "failed to vote: %s", err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func httpErr(w *http.ResponseWriter, code int, s string, args ...interface{}) {
	msg := fmt.Errorf(s, args...)
	(*w).WriteHeader(code)
	fmt.Fprintln(*w, msg)
	if *viewerRequestLogging {
		fmt.Println(fmt.Sprintf("%d: %s", code, msg))
	}
}
