module github.com/gkstretton/dark/services/dslrcapture

go 1.19

require (
	github.com/gkstretton/asol-protos v0.0.9-0.20230708095357-d91016e3f669
	github.com/gkstretton/dark/services/goo v0.0.0-20230707193130-3804eccf6636
	google.golang.org/protobuf v1.31.0
)

require (
	github.com/andreykaipov/goobs v0.12.0 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/eclipse/paho.mqtt.golang v1.4.3 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/hashicorp/logutils v1.0.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/nu7hatch/gouuid v0.0.0-20131221200532-179d4d0c4d8d // indirect
	golang.org/x/net v0.12.0 // indirect
	golang.org/x/sync v0.3.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)

// /goo is mounted by docker compose
replace github.com/gkstretton/dark/services/goo => ../goo
