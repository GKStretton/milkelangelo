package livecapture

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/gkstretton/dark/services/goo/config"
	"github.com/gkstretton/dark/services/goo/filesystem"
)

type webcamRecorder struct {
	sessionId uint64
	name      string
	cmd       *exec.Cmd
	stdin     io.WriteCloser
}

func startWebcamRecording(rtspPath string, sessionId uint64) (*webcamRecorder, error) {
	url := fmt.Sprintf("%s:8554/%s", *config.RtspHost, rtspPath)

	dir := filesystem.GetRawVideoDir(sessionId, rtspPath)
	filePath := filesystem.GetIncrementalFileName(dir, "mp4")

	fmt.Printf("Calling capture-rtsp.sh with '%s' and '%s'\n", url, filePath)

	wr := &webcamRecorder{
		sessionId: sessionId,
		name:      rtspPath,
	}
	wr.cmd = exec.Command("./scripts/capture-rtsp.sh", url, filePath)

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
	prefix := fmt.Sprintf("[Webcam Recorder] (%d - %s): ", wr.sessionId, wr.name)
	fmt.Printf(prefix+s, args...)
}
