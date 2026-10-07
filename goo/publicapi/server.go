package publicapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"runtime/debug"
	"time"
)

const (
	maxBodyBytes = 1 << 10
	claimHeader  = "X-Claim-Token"
)

type Options struct {
	Addr string
	// UiDir is the built remote control page, served at / if set.
	UiDir string
	// MediamtxURL is MediaMTX's WebRTC server. If set, signalling for the
	// cropped cameras is proxied at /video/.
	MediamtxURL string
	// CloudflareTurnKeyID and CloudflareTurnAPIToken, if both set, let
	// viewers who can't reach MediaMTX directly relay video through
	// Cloudflare's TURN service. See /api/ice-servers.
	CloudflareTurnKeyID    string
	CloudflareTurnAPIToken string
}

// Start serves the api under /api/, plus the page and video signalling when
// configured. It uses its own server and mux so that nothing registered on
// http.DefaultServeMux is exposed.
func (a *Api) Start(o Options) error {
	h, err := a.rootHandler(o)
	if err != nil {
		return err
	}

	server := &http.Server{
		Addr:              o.Addr,
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,
		// /events clears these for its stream; upgraded /video/ connections
		// have them cleared by net/http
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		l.Printf("listening on %s\n", o.Addr)
		err := server.ListenAndServe()
		l.Printf("server stopped: %v\n", err)
	}()
	return nil
}

func (a *Api) rootHandler(o Options) (http.Handler, error) {
	if o.CloudflareTurnKeyID != "" && o.CloudflareTurnAPIToken != "" {
		a.turn = newCloudflareTurn(o.CloudflareTurnKeyID, o.CloudflareTurnAPIToken)
	}

	mux := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", a.handler()))
	if o.UiDir != "" {
		mux.Handle("/", http.FileServer(http.Dir(o.UiDir)))
	}
	if o.MediamtxURL != "" {
		u, err := url.Parse(o.MediamtxURL)
		if err != nil {
			return nil, fmt.Errorf("invalid mediamtx url: %w", err)
		}
		mux.Handle("/video/", videoProxy(u))
	}
	return mux, nil
}

func (a *Api) handler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/state", method(http.MethodGet, a.getState))
	mux.Handle("/events", method(http.MethodGet, a.streamEvents))
	mux.Handle("/claim", method(http.MethodPut, a.handleClaim))
	mux.Handle("/unclaim", method(http.MethodPut, a.handleUnclaim))
	// not /collect, which ad and privacy blockers block as an analytics endpoint
	mux.Handle("/collection", method(http.MethodPost, a.handleCollect))
	mux.Handle("/dispense", method(http.MethodPost, a.handleDispense))
	mux.Handle("/goto", method(http.MethodPut, a.handleGoTo))
	mux.Handle("/ice-servers", method(http.MethodGet, a.getIceServers))

	return recoverer(cors(limitBody(mux)))
}

func method(m string, h http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != m {
			writeError(w, http.StatusMethodNotAllowed, fmt.Errorf("method %s not allowed", r.Method))
			return
		}
		h(w, r)
	})
}

// recoverer stops a panic in a handler from taking down goo.
func recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				l.Printf("panic handling %s %s: %v\n%s", r.Method, r.URL.Path, err, debug.Stack())
				writeError(w, http.StatusInternalServerError, fmt.Errorf("internal error"))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// cors allows any origin, so a page on another origin can call the api.
// Control is gated by the claim token, not cookies.
func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("Access-Control-Allow-Origin", "*")
		h.Set("Access-Control-Allow-Headers", "Content-Type, "+claimHeader)
		h.Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func limitBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
		next.ServeHTTP(w, r)
	})
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"message": err.Error()})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
