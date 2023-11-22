package mqtt

import (
	"fmt"
	"time"
)

func SubscribeBlocking(topic string, timeout time.Duration) ([]byte, error) {
	// Create a channel to receive the message payload
	payloadChan := make(chan []byte, 1)

	// Subscribe to the topic
	Subscribe(topic, func(topic string, pl []byte) {
		payloadChan <- pl
	})
	defer Unsubscribe(topic)

	// Wait for the message or timeout
	select {
	case payload := <-payloadChan:
		return payload, nil
	case <-time.After(timeout):
		return nil, fmt.Errorf("SubscribeBlocking timed out after %v", timeout)
	}
}
