package publicapi

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gkstretton/dark/services/goo/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func do(h http.Handler, method, path, token string, body any) *httptest.ResponseRecorder {
	var b []byte
	if body != nil {
		b, _ = json.Marshal(body)
	}
	req := httptest.NewRequest(method, path, bytes.NewReader(b))
	if token != "" {
		req.Header.Set(claimHeader, token)
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	return w
}

func claimToken(t *testing.T, h http.Handler, token string) string {
	w := do(h, http.MethodPut, "/claim", token, nil)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	var resp struct{ Token string }
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.NotEmpty(t, resp.Token)
	return resp.Token
}

func getState(t *testing.T, h http.Handler) state {
	var s state
	require.NoError(t, json.Unmarshal(do(h, http.MethodGet, "/state", "", nil).Body.Bytes(), &s))
	return s
}

func TestClaim(t *testing.T) {
	h := New().handler()
	assert.False(t, getState(t, h).Claimed)

	alice := claimToken(t, h, "")
	assert.True(t, getState(t, h).Claimed)

	// renewing keeps the same token
	assert.Equal(t, alice, claimToken(t, h, alice))

	// others can't claim, unclaim or command
	assert.Equal(t, http.StatusForbidden, do(h, http.MethodPut, "/claim", "", nil).Code)
	assert.Equal(t, http.StatusForbidden, do(h, http.MethodPut, "/claim", "guess", nil).Code)
	assert.Equal(t, http.StatusForbidden, do(h, http.MethodPut, "/unclaim", "guess", nil).Code)
	assert.Equal(t, http.StatusForbidden, do(h, http.MethodPost, "/collection", "", map[string]int{"id": 2}).Code)

	assert.Equal(t, http.StatusAccepted, do(h, http.MethodPut, "/unclaim", alice, nil).Code)
	assert.False(t, getState(t, h).Claimed)

	// a released token no longer works
	assert.Equal(t, http.StatusForbidden, do(h, http.MethodPost, "/collection", alice, map[string]int{"id": 2}).Code)
	bob := claimToken(t, h, "")
	assert.NotEqual(t, alice, bob)
}

func TestCommandRejectedByControl(t *testing.T) {
	h := New().handler()
	token := claimToken(t, h, "")

	assert.Equal(t, http.StatusBadRequest, do(h, http.MethodPost, "/collection", token, map[string]int{}).Code)

	// control is disabled by default, so the command is rejected with a reason
	w := do(h, http.MethodPost, "/collection", token, map[string]int{"id": 2})
	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Contains(t, w.Body.String(), "disabled")
}

func TestMethodNotAllowed(t *testing.T) {
	h := New().handler()
	assert.Equal(t, http.StatusMethodNotAllowed, do(h, http.MethodGet, "/claim", "", nil).Code)
	assert.Equal(t, http.StatusNoContent, do(h, http.MethodOptions, "/claim", "", nil).Code)
}

func TestEvents(t *testing.T) {
	a := New()
	srv := httptest.NewServer(a.handler())
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/events")
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, "text/event-stream", resp.Header.Get("Content-Type"))

	events := make(chan state)
	go func() {
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			if data, ok := strings.CutPrefix(scanner.Text(), "data: "); ok {
				var s state
				if json.Unmarshal([]byte(data), &s) == nil {
					events <- s
				}
			}
		}
	}()

	next := func() state {
		select {
		case s := <-events:
			return s
		case <-time.After(2 * time.Second):
			t.Fatal("timed out waiting for event")
			return state{}
		}
	}

	// initial state on connect
	assert.Equal(t, float32(0), next().GooState.X)

	a.UpdateState(func(s *types.GooState) { s.X = 0.5 })
	assert.Equal(t, float32(0.5), next().GooState.X)

	_, err = a.claim("")
	require.NoError(t, err)
	assert.True(t, next().Claimed)
}

func TestRootHandler(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "index.html"), []byte("page"), 0o644))
	h, err := New().rootHandler(Options{UiDir: dir})
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, do(h, http.MethodGet, "/api/state", "", nil).Code)
	assert.Equal(t, http.StatusNotFound, do(h, http.MethodGet, "/state", "", nil).Code)

	w := do(h, http.MethodGet, "/", "", nil)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "page", w.Body.String())

	// without a ui dir, only the api is served
	apiOnly, err := New().rootHandler(Options{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, do(apiOnly, http.MethodGet, "/", "", nil).Code)
	assert.Equal(t, http.StatusNotFound, do(apiOnly, http.MethodGet, "/video/top-cam-crop/ws", "", nil).Code)
}
