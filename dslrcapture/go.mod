module github.com/gkstretton/dark/services/dslrcapture

go 1.21

require (
	github.com/gkstretton/asol-protos v0.0.9-0.20231121080309-bd93879457b6
	github.com/gkstretton/dark/services/goo v0.0.0-20231114224855-2d1a2074d446
	google.golang.org/protobuf v1.31.0
)

require (
	github.com/eclipse/paho.mqtt.golang v1.4.3 // indirect
	github.com/gorilla/websocket v1.5.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mailjet/mailjet-apiv3-go/v3 v3.2.0 // indirect
	github.com/mailjet/mailjet-apiv3-go/v4 v4.0.1 // indirect
	golang.org/x/net v0.19.0 // indirect
	golang.org/x/sync v0.5.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)

// /goo is mounted by docker compose
replace github.com/gkstretton/dark/services/goo => ../goo

replace github.com/gkstretton/asol-protos => ../asol-protos
