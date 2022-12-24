package livecapture

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/gkstretton/dark/services/goo/config"
)

type webcamRecorder struct {
	sessionId uint64
	name      string
	cmd       *exec.Cmd
	stdin     io.WriteCloser
}

func startWebcamRecording(rtspPath string, sessionId uint64) (*webcamRecorder, error) {
	url := fmt.Sprintf("%s:8554/%s", *config.RtspHost, rtspPath)
	outDir := fmt.Sprintf("%s/%d/%s/%s/",
		*config.SessionBasePath,
		uint64(sessionId),
		*config.RawFootagePath,
		rtspPath,
	)

	//todo: incremental file name for resumes
	fileName := "1.mp4"

	fmt.Printf("Calling capture-rtsp.sh with '%s' and '%s'\n", url, outDir+fileName)

	wr := &webcamRecorder{
		sessionId: sessionId,
		name:      rtspPath,
	}
	wr.cmd = exec.Command("./scripts/capture-rtsp.sh", url, outDir+fileName)

	wr.cmd.Stdout = os.Stdout
	wr.cmd.Stderr = os.Stderr

	var err error
	wr.stdin, err = wr.cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stdin pipe in startWebcamRecording: %v", err)
	}

	err = wr.cmd.Start()
	if err != nil {
		return nil, err
	}

	wr.log("capture started")

	return wr, nil
}

func (wr *webcamRecorder) Stop() {
	_, err := wr.stdin.Write([]byte("q"))
	if err != nil {
		wr.log("failed to terminate video recording by stdin: %v", err)
		wr.cmd.Process.Kill()
		return
	}
	err = wr.cmd.Wait()
	if err != nil {
		wr.log("failed to gracefully stop recording: %v", err)
		wr.cmd.Process.Kill()
		return
	}
	wr.log("gracefully stopped recording")
}

func (wr *webcamRecorder) log(s string, args ...interface{}) {
	msg := fmt.Sprintf("[Webcam Recorder] (%d - %s): "+s+"\n", wr.sessionId, wr.name, args)
	fmt.Println(msg)
}
