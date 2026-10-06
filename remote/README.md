# Milkelangelo remote control

A web page for controlling the robot remotely, served by goo along with its
API (`goo/publicapi`, spec in `openapi/openapi.yaml`).

- `frontend/` - the control page (React + Vite)
- `openapi/` - the API spec

One person controls the robot at a time: "Take control" claims it, and the
page renews the claim every few seconds until they give it up or close the
page. Commands are checked by `goo/control`, so the page only shows what's
allowed; the robot enforces it.

The API does no authentication. Whoever can reach it can use it, so put it
behind something that controls access before exposing it beyond your own
networks. Run locally, goo only listens on `127.0.0.1:8789`; deployed with
docker-compose it listens on all interfaces, like the interface.

## Video

The page plays the cropped top camera (`top-cam-crop`) as the bowl view, and
the cropped front camera (`front-cam-crop`) beside it, from MediaMTX over
WebRTC. Each WebRTC handshake (a websocket) goes through goo at
`/video/<stream>/ws` on the page's own origin, so it works wherever the page
is reachable, including over https. goo proxies only those two streams, to
`MEDIAMTX_URL`. Set `VITE_WEBRTC_URL` at build time to bypass goo and talk to
MediaMTX directly.

The video itself doesn't go through goo. Browser and MediaMTX connect
directly, finding a route with the STUN server in MediaMTX's
`webrtcICEServers`. This is the same route the cloud interface's video takes. It
works on most networks; ones that block it (strict corporate or mobile
networks) would need a TURN relay added to `webrtcICEServers`. If video can't
connect, the page shows a plain bowl with the pipette position.

## Developing

0. Stop goo on milkelangelo to avoid interference
1. Run goo: `make goo` (or `cd goo && go run .`)
2. Run the page: `make remote`, then open `localhost:3000`. The dev server
   proxies `/api` to goo at `127.0.0.1:8789` (override with `GOO_API`).
3. Wake the machine and press "Enable Remote Control" in the interface.

To try goo serving the built page itself:

```bash
cd remote/frontend && npm run build
cd goo && PUBLIC_UI_DIR=../remote/frontend/dist go run .
# open http://localhost:8789
```

The API can also be used directly:

```bash
api=http://localhost:8789/api
curl -N $api/events &                      # follow state
token=$(curl -s -X PUT $api/claim | jq -r .token)
curl -X POST $api/collection -H "X-Claim-Token: $token" -d '{"id":2}'
```

## Deploying

goo's Docker image builds the page and serves it from `/app/remote`, at
`http://milkelangelo:8789` from the LAN or tailscale. Set
`PUBLIC_API_ADDR=127.0.0.1:8789` in `.env` to keep it to that machine, or
`ENABLE_PUBLIC_API=false` to turn it off.

## goo configuration

| Env var | Default | |
|---|---|---|
| `ENABLE_PUBLIC_API` | `true` | serve the page and API |
| `PUBLIC_API_ADDR` | `127.0.0.1:8789` (`:8789` in docker-compose) | listen address |
| `PUBLIC_UI_DIR` | none (`/app/remote` in docker) | built page to serve at `/` |
| `MEDIAMTX_URL` | `http://milkelangelo:8889` (`http://127.0.0.1:8889` in docker) | MediaMTX WebRTC server for video. Empty disables video |
