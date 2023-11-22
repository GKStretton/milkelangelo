package livecapture

import (
	"fmt"
	"sync"

	"github.com/gkstretton/dark/services/goo/config"
	"github.com/gkstretton/dark/services/goo/session"
)

func (r *recorder) record(id session.ID) {
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

	<-r.stopRecording
}
