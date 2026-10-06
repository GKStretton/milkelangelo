package publicapi

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"slices"
	"strings"
)

// videoStreams are the MediaMTX streams whose WebRTC signalling is exposed:
// the cropped bowl (top) and front cameras. The uncropped ones stay on the LAN.
var videoStreams = []string{"top-cam-crop", "front-cam-crop"}

// videoProxy forwards WebRTC signalling websockets at /video/<stream>/ws to
// MediaMTX, so the page can reach it on its own origin. The media itself
// doesn't pass through here: it goes between the browser and MediaMTX
// directly, using the ICE servers MediaMTX is configured with.
//
// The websocket stays open for the whole viewing session. net/http clears
// the server's timeouts when the proxy hijacks the connection to upgrade it.
func videoProxy(mediamtx *url.URL) http.Handler {
	proxy := &httputil.ReverseProxy{
		Rewrite: func(pr *httputil.ProxyRequest) {
			pr.SetURL(mediamtx)
			pr.Out.Host = mediamtx.Host
		},
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		stream, ok := strings.CutSuffix(strings.TrimPrefix(r.URL.Path, "/video/"), "/ws")
		if !ok || !slices.Contains(videoStreams, stream) {
			http.NotFound(w, r)
			return
		}

		r.URL.Path = "/" + stream + "/ws"
		r.URL.RawPath = ""
		proxy.ServeHTTP(w, r)
	})
}
