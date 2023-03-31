package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"sync"
	"time"

	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/asol-protos/go/topics_backend"
	"github.com/gkstretton/dark/services/dslrcapture/mqtt"
)

var (
	isRecording     bool
	stopRecording   chan bool
	mutex           *sync.Mutex
	captureInterval = flag.Int("captureInterval", 10, "The interval in seconds between captures")
)

func main() {
	flag.Parse()
	mqtt.Start()

	mqtt.Subscribe(topics_backend.TOPIC_SESSION_STATUS_RESP_RAW, func(topic string, payload []byte) {
		var status *machinepb.SessionStatus
		err := json.Unmarshal(payload, &status)
		if err != nil {
			fmt.Printf("Error unmarshalling session status response: %v\n", err)
			return
		}
		handleSessionStatus(status)
	})

	// special handler for the crop config preview capture
	registerDslrPreviewHandler()

	mqtt.Publish(topics_backend.TOPIC_SESSION_STATUS_GET, "")
}

func handleSessionStatus(status *machinepb.SessionStatus) {
	fmt.Printf("handling session status: %v\n", status)
	mutex.Lock()
	defer mutex.Unlock()

	if isRecording && (status.Complete || status.Paused) {
		<-stopRecording
	}

	if !isRecording && !status.Complete && !status.Paused {
		go captureLoop(status.Id)
	}
}

func captureLoop(sessionNumber uint64) {
	setIsRecording(true)
	defer setIsRecording(false)

	next := time.After(time.Millisecond)
	for {
		next = time.After(time.Duration(*captureInterval) * time.Second)
		select {
		case <-stopRecording:
			return
		case <-next:
			captureSessionImage(sessionNumber)
		}
	}
}

func setIsRecording(b bool) {
	mutex.Lock()
	defer mutex.Unlock()
	isRecording = b
}
