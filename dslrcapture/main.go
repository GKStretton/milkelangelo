package main

import (
	"flag"
	"fmt"
	"sync"
	"time"

	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/asol-protos/go/topics_backend"
	"github.com/gkstretton/dark/services/goo/mqtt"
	"google.golang.org/protobuf/proto"
)

var (
	isRecording     bool
	stopRecording   chan bool = make(chan bool)
	mutex           sync.Mutex
	captureInterval = flag.Int("captureInterval", 30, "The interval in seconds between captures")
	host            = flag.String("host", "localhost", "the hostname of the mqtt broker")
)

func main() {
	flag.Parse()
	mqtt.Start(*host)

	mqtt.Subscribe(topics_backend.TOPIC_SESSION_STATUS_RESP_RAW, func(topic string, payload []byte) {
		status := &machinepb.SessionStatus{}
		err := proto.Unmarshal(payload, status)
		if err != nil {
			fmt.Printf("Error unmarshalling session status response: %v\n", err)
			return
		}
		handleSessionStatus(status)
	})

	// special handler for the crop config preview capture
	registerDslrPreviewHandler()

	time.Sleep(time.Second)
	mqtt.Publish(topics_backend.TOPIC_SESSION_STATUS_GET, "")

	for {
		time.Sleep(time.Second)
	}
}

func handleSessionStatus(status *machinepb.SessionStatus) {
	fmt.Printf("handling session status: %v\n", status)
	mutex.Lock()
	defer mutex.Unlock()

	if isRecording && (status.Complete || status.Paused) {
		fmt.Printf("issuing stop recording\n")
		stopRecording <- true
		fmt.Printf("stop recording passed\n")
		return
	}

	if !isRecording && !status.Complete && !status.Paused {
		fmt.Printf("launching capture loop\n")
		go captureLoop(status.Id)
	}
}

func captureLoop(sessionNumber uint64) {
	setIsRecording(true)
	defer setIsRecording(false)

	next := time.After(time.Millisecond)
	for {
		select {
		case <-stopRecording:
			fmt.Printf("stopping recording, returning\n")
			return
		case <-next:
			next = time.After(time.Duration(*captureInterval) * time.Second)
			captureSessionImage(sessionNumber)
		}
	}
}

func setIsRecording(b bool) {
	mutex.Lock()
	defer mutex.Unlock()
	isRecording = b
	fmt.Printf("set recording to %t\n", b)
}
