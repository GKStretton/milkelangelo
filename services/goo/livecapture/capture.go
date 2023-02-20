package livecapture

import (
	"fmt"
	"sync"
	"time"

	"github.com/gkstretton/dark/services/goo/config"
	"github.com/gkstretton/dark/services/goo/session"
)

func (r *recorder) record(id session.ID) {
	// after testing, it seems that this isn't needed:
	// setDslrState(true)
	// defer setDslrState(false)

	r.isRecording = true
	defer func() { r.isRecording = false }()

	// webcam recordings

	var topRecording, frontRecording *webcamRecorder
	var err error
	wg := sync.WaitGroup{}
	go func() {
		defer wg.Done()
		topRecording, err = startWebcamRecording(*config.TopCamRtspPath, uint64(id))
		if err != nil {
			fmt.Printf("failed to start top webcam recording: %v\n", err)
		}
	}()
	go func() {
		defer wg.Done()
		frontRecording, err = startWebcamRecording(*config.FrontCamRtspPath, uint64(id))
		if err != nil {
			fmt.Printf("failed to start front webcam recording: %v\n", err)
		}
	}()
	wg.Wait()

	defer func() {
		go topRecording.Stop()
		go frontRecording.Stop()
	}()

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
