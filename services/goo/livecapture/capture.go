package livecapture

import (
	"fmt"
	"time"

	"github.com/gkstretton/dark/services/goo/config"
	"github.com/gkstretton/dark/services/goo/session"
)

func (r *recorder) record(id session.ID) {
	setDslrState(true)
	defer setDslrState(false)

	r.isRecording = true
	defer func() { r.isRecording = false }()

	// webcam recordings

	topRecording, err := startWebcamRecording(*config.TopCamRtspPath, uint64(id))
	if err != nil {
		fmt.Printf("failed to start top webcam recording: %v\n", err)
	}
	frontRecording, err := startWebcamRecording(*config.FrontCamRtspPath, uint64(id))
	if err != nil {
		fmt.Printf("failed to start front webcam recording: %v\n", err)
	}

	defer topRecording.Stop()
	defer frontRecording.Stop()

	// Regular image capture

	next := time.After(0)
	for {
		select {
		case <-next:
			next = time.After(time.Second * time.Duration(*captureInterval))
			captureSessionImage(uint64(id))
		case <-r.stopRecording:
			return
		}
	}
}
