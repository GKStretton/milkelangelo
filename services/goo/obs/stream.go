package obs

import (
	"fmt"

	"github.com/andreykaipov/goobs/api/requests/stream"
	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/asol-protos/go/topics_backend"
	"github.com/gkstretton/dark/services/goo/mqtt"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func startStream(topic string, payload []byte) {
	if c == nil {
		fmt.Println("obs client is nil")
		return
	}

	handleSessionEvent(sm)
	setCropConfig()

	_, err := c.Stream.StartStream(&stream.StartStreamParams{})
	if err != nil {
		fmt.Printf("failed to start streaming: %v\n", err)
		return
	}
	fmt.Printf("sent start streaming request\n")

	onStreamStart()
}

func endStream(topic string, payload []byte) {
	if c == nil {
		fmt.Println("obs client is nil")
		return
	}

	_, err := c.Stream.StopStream()
	if err != nil {
		fmt.Printf("failed to stop streaming: %v\n", err)
		return
	}
	fmt.Printf("sent stop streaming request\n")

	onStreamEnd()
}

func isStreamLive() bool {
	if c == nil {
		fmt.Println("obs client is nil")
		return false
	}

	resp, err := c.Stream.GetStreamStatus()
	if err != nil {
		fmt.Printf("failed to get stream status: %v\n", err)
		return false
	}

	return resp.OutputActive
}

func publishStreamStatus() {
	s := &machinepb.StreamStatus{
		Live: isStreamLive(),
	}

	// protobuf
	b, err := proto.Marshal(s)
	if err != nil {
		fmt.Printf("error marshalling stream status as protobuf: %v\n", err)
	}
	if err = mqtt.Publish(topics_backend.TOPIC_STREAM_STATUS_RESP_RAW, b); err != nil {
		fmt.Printf("error publishing stream status: %v\n", err)
	}

	// json
	m := protojson.MarshalOptions{
		Multiline:       true,
		UseProtoNames:   true,
		Indent:          "\t",
		EmitUnpopulated: true,
	}
	j, err := m.Marshal(s)
	if err != nil {
		fmt.Printf("error marshalling stream status to json: %v\n", err)
	}
	if err = mqtt.Publish(topics_backend.TOPIC_STREAM_STATUS_RESP_JSON, j); err != nil {
		fmt.Printf("error publishing stream status: %v\n", err)
	}
}
