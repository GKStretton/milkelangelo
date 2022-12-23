package livecapture

import (
	"fmt"
	"time"

	"github.com/gkstretton/dark/services/goo/session"
)

func (r *recorder) record(id session.ID) {
	defer stopVideoRecording()
	defer func() { r.isRecording = false }()

	r.isRecording = true
	startVideoRecording()

	next := time.After(0)
	for {
		select {
		case <-next:
			next = time.After(time.Second * time.Duration(*captureInterval))
			captureImage()
		case <-r.stopRecording:
			return
		}
	}
}

func captureImage() {
	fmt.Println("capture image (not implemented)")
}

func startVideoRecording() {
	fmt.Println("start recording (not implemented)")
}

func stopVideoRecording() {
	fmt.Println("stop recording (not implemented)")
}
