package obs

import (
	"fmt"
)

func startStream(topic string, payload []byte) {
	if c == nil {
		fmt.Println("obs client is nil")
		return
	}
	handleSessionEvent(sm)
	setCropConfig()

	// _, err := c.Streaming.StartStreaming(&streaming.StartStreamingParams{})
	// if err != nil {
	// 	fmt.Printf("failed to start streaming: %v\n", err)
	// 	return
	// }
	fmt.Printf("sent start streaming request\n")
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
}
