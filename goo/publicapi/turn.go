package publicapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	cloudflareTurnAPI = "https://rtc.live.cloudflare.com/v1/turn/keys"
	// how long generated credentials last
	turnCredentialTTL = 24 * time.Hour
	// generate new credentials when the cached ones have less than this left,
	// so a viewer never receives nearly-expired ones
	turnRefreshMargin = 2 * time.Hour
)

var errNoTurn = errors.New("no TURN relay configured")

// cloudflareTurn generates short-lived Cloudflare TURN credentials, so that
// video can be relayed for viewers who can't reach MediaMTX directly. They
// are cached and shared between viewers, so Cloudflare is called about once
// a day at most.
type cloudflareTurn struct {
	keyID    string
	apiToken string
	apiBase  string
	client   *http.Client

	lock    sync.Mutex
	cached  json.RawMessage
	expires time.Time
}

func newCloudflareTurn(keyID, apiToken string) *cloudflareTurn {
	return &cloudflareTurn{
		keyID:    keyID,
		apiToken: apiToken,
		apiBase:  cloudflareTurnAPI,
		client:   &http.Client{Timeout: 10 * time.Second},
	}
}

// iceServers returns an RTCPeerConnection iceServers list.
func (t *cloudflareTurn) iceServers(ctx context.Context) (json.RawMessage, error) {
	if t == nil {
		return nil, errNoTurn
	}

	t.lock.Lock()
	defer t.lock.Unlock()

	if t.cached != nil && time.Until(t.expires) > turnRefreshMargin {
		return t.cached, nil
	}

	servers, err := t.generate(ctx)
	if err != nil {
		return nil, err
	}
	t.cached = servers
	t.expires = time.Now().Add(turnCredentialTTL)
	return servers, nil
}

func (t *cloudflareTurn) generate(ctx context.Context) (json.RawMessage, error) {
	body, _ := json.Marshal(map[string]int{"ttl": int(turnCredentialTTL.Seconds())})
	url := fmt.Sprintf("%s/%s/credentials/generate-ice-servers", t.apiBase, t.keyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+t.apiToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("requesting turn credentials: %w", err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(io.LimitReader(resp.Body, 64<<10))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("turn credentials: unexpected status %d: %s", resp.StatusCode, b)
	}

	var parsed struct {
		IceServers []iceServer `json:"iceServers"`
	}
	if err := json.Unmarshal(b, &parsed); err != nil {
		return nil, fmt.Errorf("parsing turn credentials: %w", err)
	}
	if len(parsed.IceServers) == 0 {
		return nil, errors.New("turn credentials: no ice servers returned")
	}

	for i := range parsed.IceServers {
		parsed.IceServers[i].URLs = withoutPort53(parsed.IceServers[i].URLs)
	}
	return json.Marshal(parsed.IceServers)
}

type iceServer struct {
	URLs       []string `json:"urls"`
	Username   string   `json:"username,omitempty"`
	Credential string   `json:"credential,omitempty"`
}

// withoutPort53 drops urls on port 53, which browsers refuse to use.
func withoutPort53(urls []string) []string {
	kept := urls[:0]
	for _, u := range urls {
		hostPort, _, _ := strings.Cut(u, "?")
		if strings.HasSuffix(hostPort, ":53") {
			continue
		}
		kept = append(kept, u)
	}
	return kept
}

// getIceServers serves {"iceServers": [...]} for the page to use for video,
// or 404 if no relay is configured.
func (a *Api) getIceServers(w http.ResponseWriter, r *http.Request) {
	servers, err := a.turn.iceServers(r.Context())
	if errors.Is(err, errNoTurn) {
		writeError(w, http.StatusNotFound, err)
		return
	}
	if err != nil {
		l.Printf("failed to get turn credentials: %v\n", err)
		writeError(w, http.StatusBadGateway, errors.New("couldn't get relay credentials"))
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, http.StatusOK, map[string]json.RawMessage{"iceServers": servers})
}
