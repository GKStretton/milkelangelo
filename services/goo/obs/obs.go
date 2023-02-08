package obs

import (
	"fmt"

	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/streaming"
	"github.com/andreykaipov/goobs/api/typedefs"
	"github.com/gkstretton/dark/services/goo/config"
	"github.com/gkstretton/dark/services/goo/mqtt"
	"github.com/gkstretton/dark/services/goo/session"
)

var c *goobs.Client

func Run(sm *session.SessionManager) {
	fmt.Println("Running OBS controller")

	var err error
	c, err = goobs.New("localhost:4444")
	if err != nil {
		fmt.Printf("failed to create obs ws client: %v\n", err)
		return
	}
	resp, err := c.General.GetVersion()
	if err != nil {
		fmt.Printf("failed to get version: %v\n", err)
		return
	}
	fmt.Printf("OBS version: %s\n", resp.ObsStudioVersion)
	fmt.Printf("OBSws version: %s\n", resp.ObsWebsocketVersion)

	mqtt.Subscribe(config.TOPIC_STREAM_START, startStream)
	mqtt.Subscribe(config.TOPIC_STREAM_END, endStream)
}

func startStream(topic string, payload []byte) {
	_, err := c.Streaming.StartStreaming(&streaming.StartStreamingParams{
		Stream: &streaming.Stream{
			Metadata: map[string]interface{}{},
			Settings: &typedefs.StreamSettings{},
		},
	})
	if err != nil {
		fmt.Printf("failed to start streaming: %v\n", err)
	}
}

func endStream(topic string, payload []byte) {
	_, err := c.Streaming.StopStreaming()
	if err != nil {
		fmt.Printf("failed to stop streaming: %v\n", err)
	}
}
