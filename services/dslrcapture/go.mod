module github.com/gkstretton/dark/services/dslrcapture

go 1.19

require (
	github.com/gkstretton/asol-protos v0.0.9-0.20230521150831-de303624dcd3
	github.com/gkstretton/dark/services/goo v0.0.0-20230520063052-2b9ce7e24769
	google.golang.org/protobuf v1.30.0
)

require (
	github.com/eclipse/paho.mqtt.golang v1.4.2 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sync v0.2.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// /goo is mounted by docker compose
replace github.com/gkstretton/dark/services/goo => ../goo
