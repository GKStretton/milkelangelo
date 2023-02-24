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

	r.setIsRecording(true)
	defer r.setIsRecording(false)

	// webcam recordings

	var topRecording, frontRecording *webcamRecorder
	var errTop, errFront error
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		topRecording, errTop = startWebcamRecording(*config.TopCamRtspPath, uint64(id))
		if errTop != nil {
			fmt.Printf("failed to start top webcam recording: %v\n", errTop)
		}
	}()
	go func() {
		defer wg.Done()
		frontRecording, errFront = startWebcamRecording(*config.FrontCamRtspPath, uint64(id))
		if errFront != nil {
			fmt.Printf("failed to start front webcam recording: %v\n", errFront)
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
