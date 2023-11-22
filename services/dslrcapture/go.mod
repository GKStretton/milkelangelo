module github.com/gkstretton/dark/services/dslrcapture

go 1.19

require (
	github.com/gkstretton/asol-protos v0.0.9-0.20231121080309-bd93879457b6
	github.com/gkstretton/dark/services/goo v0.0.0-20231114224855-2d1a2074d446
	google.golang.org/protobuf v1.31.0
)

require (
	github.com/Davincible/goinsta/v3 v3.2.6 // indirect
	github.com/andreykaipov/goobs v0.12.1 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/chromedp/cdproto v0.0.0-20231114014204-3e458d5176f9 // indirect
	github.com/chromedp/chromedp v0.9.3 // indirect
	github.com/chromedp/sysutil v1.0.0 // indirect
	github.com/eclipse/paho.mqtt.golang v1.4.3 // indirect
	github.com/gempir/go-twitch-irc/v4 v4.0.0 // indirect
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/gobwas/ws v1.3.1 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/gorilla/websocket v1.5.1 // indirect
	github.com/hashicorp/logutils v1.0.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/nicklaw5/helix/v2 v2.25.2 // indirect
	github.com/nu7hatch/gouuid v0.0.0-20131221200532-179d4d0c4d8d // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/net v0.18.0 // indirect
	golang.org/x/sync v0.5.0 // indirect
	golang.org/x/sys v0.14.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)

// /goo is mounted by docker compose
replace github.com/gkstretton/dark/services/goo => ../goo
