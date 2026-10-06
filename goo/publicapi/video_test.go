package publicapi

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeMediamtx accepts a websocket-style upgrade at /<stream>/ws and echoes
// whatever it receives.
func fakeMediamtx(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/ws") || r.Header.Get("Upgrade") != "websocket" {
			http.Error(w, "unexpected "+r.URL.Path, http.StatusBadRequest)
			return
		}
		conn, buf, err := http.NewResponseController(w).Hijack()
		require.NoError(t, err)
		defer conn.Close()
		fmt.Fprint(buf, "HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\n\r\n")
		buf.Flush()
		io.Copy(conn, buf)
	}))
}

// serve runs the root handler with short timeouts, standing in for the
// real server's 10s ones.
func serve(t *testing.T, o Options) string {
	h, err := New().rootHandler(o)
	require.NoError(t, err)
	srv := httptest.NewUnstartedServer(h)
	srv.Config.ReadTimeout = 200 * time.Millisecond
	srv.Config.WriteTimeout = 200 * time.Millisecond
	srv.Start()
	t.Cleanup(srv.Close)
	return srv.Listener.Addr().String()
}

func upgrade(t *testing.T, addr, path string) (net.Conn, *bufio.Reader, int) {
	conn, err := net.Dial("tcp", addr)
	require.NoError(t, err)
	t.Cleanup(func() { conn.Close() })

	fmt.Fprintf(conn, "GET %s HTTP/1.1\r\nHost: example\r\nConnection: Upgrade\r\nUpgrade: websocket\r\n\r\n", path)
	r := bufio.NewReader(conn)
	resp, err := http.ReadResponse(r, nil)
	require.NoError(t, err)
	return conn, r, resp.StatusCode
}

func TestVideoProxyOutlivesTimeouts(t *testing.T) {
	mediamtx := fakeMediamtx(t)
	defer mediamtx.Close()
	addr := serve(t, Options{MediamtxURL: mediamtx.URL})

	conn, r, status := upgrade(t, addr, "/video/top-cam-crop/ws")
	require.Equal(t, http.StatusSwitchingProtocols, status)

	// well past the server's read and write timeouts
	time.Sleep(500 * time.Millisecond)

	_, err := conn.Write([]byte("ping\n"))
	require.NoError(t, err)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	line, err := r.ReadString('\n')
	require.NoError(t, err)
	assert.Equal(t, "ping\n", line)
}

func TestVideoProxyOnlyCroppedCameras(t *testing.T) {
	mediamtx := fakeMediamtx(t)
	defer mediamtx.Close()
	addr := serve(t, Options{MediamtxURL: mediamtx.URL})

	for _, path := range []string{"/video/top-cam-crop/ws", "/video/front-cam-crop/ws"} {
		_, _, status := upgrade(t, addr, path)
		assert.Equal(t, http.StatusSwitchingProtocols, status, path)
	}

	for _, path := range []string{
		"/video/top-cam/ws",
		"/video/front-cam/ws",
		"/video/top-cam-crop/",
		"/video/../top-cam/ws",
	} {
		_, _, status := upgrade(t, addr, path)
		assert.NotEqual(t, http.StatusSwitchingProtocols, status, path)
	}
}
