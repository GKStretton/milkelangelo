package gooapi

import "net/http"

func (g *connectedGooApi) updateHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "chunked")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	_, err := g.verifyInternalRequest(r)
	if err != nil {
		httpErr(&w, http.StatusUnauthorized, "failed to verify internal (goo) request: %v", err)
		return
	}

	// todo: unmarshal data

	w.WriteHeader(http.StatusNotImplemented)
	w.(http.Flusher).Flush()
}
