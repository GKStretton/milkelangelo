package publicapi

import (
	"fmt"
	"net/http"
	"time"
)

const keepAliveInterval = 15 * time.Second

// streamEvents sends state as server-sent events: once on connect, then
// whenever it changes.
func (a *Api) streamEvents(w http.ResponseWriter, r *http.Request) {
	rc := http.NewResponseController(w)
	// long-lived response, so lift the server's write timeout
	if err := rc.SetWriteDeadline(time.Time{}); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	c := a.subscribe()
	defer a.unsubscribe(c)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)

	keepAlive := time.NewTicker(keepAliveInterval)
	defer keepAlive.Stop()

	send := func() error {
		data, err := a.marshalState()
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w, "event: state\ndata: %s\n\n", data); err != nil {
			return err
		}
		return rc.Flush()
	}

	if err := send(); err != nil {
		return
	}
	for {
		select {
		case <-r.Context().Done():
			return
		case <-c:
			if err := send(); err != nil {
				return
			}
		case <-keepAlive.C:
			if _, err := fmt.Fprint(w, ": keep-alive\n\n"); err != nil {
				return
			}
			if err := rc.Flush(); err != nil {
				return
			}
		}
	}
}
