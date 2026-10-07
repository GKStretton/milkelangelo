package publicapi

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const cloudflareResponse = `{"iceServers":[{"urls":["stun:stun.cloudflare.com:3478","turn:turn.cloudflare.com:3478?transport=udp","turn:turn.cloudflare.com:53?transport=udp","turns:turn.cloudflare.com:5349?transport=tcp"],"username":"u","credential":"c"}]}`

// fakeCloudflare counts credential requests and checks they're well formed.
func fakeCloudflare(t *testing.T, status int) (*httptest.Server, *atomic.Int32) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		assert.Equal(t, "/key-id/credentials/generate-ice-servers", r.URL.Path)
		assert.Equal(t, "Bearer api-token", r.Header.Get("Authorization"))
		body, _ := io.ReadAll(r.Body)
		assert.JSONEq(t, `{"ttl":86400}`, string(body))

		w.WriteHeader(status)
		if status == http.StatusCreated {
			io.WriteString(w, cloudflareResponse)
		}
	}))
	t.Cleanup(srv.Close)
	return srv, &calls
}

func apiWithTurn(t *testing.T, cloudflare *httptest.Server) http.Handler {
	a := New()
	h, err := a.rootHandler(Options{CloudflareTurnKeyID: "key-id", CloudflareTurnAPIToken: "api-token"})
	require.NoError(t, err)
	a.turn.apiBase = cloudflare.URL
	return h
}

func TestIceServers(t *testing.T) {
	cloudflare, calls := fakeCloudflare(t, http.StatusCreated)
	h := apiWithTurn(t, cloudflare)

	w := do(h, http.MethodGet, "/api/ice-servers", "", nil)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	var resp struct{ IceServers []iceServer }
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Len(t, resp.IceServers, 1)
	assert.Equal(t, "u", resp.IceServers[0].Username)
	assert.Equal(t, "c", resp.IceServers[0].Credential)
	// port 53 is dropped, browsers won't use it
	assert.Equal(t, []string{
		"stun:stun.cloudflare.com:3478",
		"turn:turn.cloudflare.com:3478?transport=udp",
		"turns:turn.cloudflare.com:5349?transport=tcp",
	}, resp.IceServers[0].URLs)

	// cached for later viewers
	require.Equal(t, http.StatusOK, do(h, http.MethodGet, "/api/ice-servers", "", nil).Code)
	assert.Equal(t, int32(1), calls.Load())
}

func TestIceServersRefreshedBeforeExpiry(t *testing.T) {
	cloudflare, calls := fakeCloudflare(t, http.StatusCreated)
	a := New()
	_, err := a.rootHandler(Options{CloudflareTurnKeyID: "key-id", CloudflareTurnAPIToken: "api-token"})
	require.NoError(t, err)
	a.turn.apiBase = cloudflare.URL

	_, err = a.turn.iceServers(context.Background())
	require.NoError(t, err)

	// nearly expired: generate fresh ones
	a.turn.expires = time.Now().Add(turnRefreshMargin - time.Minute)
	_, err = a.turn.iceServers(context.Background())
	require.NoError(t, err)
	assert.Equal(t, int32(2), calls.Load())
}

func TestIceServersCloudflareError(t *testing.T) {
	cloudflare, _ := fakeCloudflare(t, http.StatusUnauthorized)
	h := apiWithTurn(t, cloudflare)

	w := do(h, http.MethodGet, "/api/ice-servers", "", nil)
	assert.Equal(t, http.StatusBadGateway, w.Code)
	// the api token isn't leaked in the error
	assert.NotContains(t, w.Body.String(), "api-token")
}

func TestIceServersNotConfigured(t *testing.T) {
	h, err := New().rootHandler(Options{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, do(h, http.MethodGet, "/api/ice-servers", "", nil).Code)
}
