package livecapture

import (
	"fmt"
	"io"
	"os/exec"
	"time"

	"github.com/gkstretton/dark/services/goo/config"
	"github.com/gkstretton/dark/services/goo/filesystem"
)

type webcamRecorder struct {
	sessionId uint64
	name      string
	cmd       *exec.Cmd
	stdin     io.WriteCloser
	filePath  string
}

func startWebcamRecording(rtspPath string, sessionId uint64) (*webcamRecorder, error) {
	url := fmt.Sprintf("%s:8554/%s", *config.RtspHost, rtspPath)

	dir := filesystem.GetRawVideoDir(sessionId, rtspPath)
	filePath := filesystem.GetIncrementalFileName(dir, "mp4")

	fmt.Printf("Calling capture-rtsp.sh with '%s' and '%s'\n", url, filePath)

	wr := &webcamRecorder{
		sessionId: sessionId,
		name:      rtspPath,
		filePath:  filePath,
	}
	wr.cmd = exec.Command("./scripts/capture-rtsp.sh", url, filePath)

	// wr.cmd.Stdout = os.Stdout
	// wr.cmd.Stderr = os.Stderr

	var err error
	wr.stdin, err = wr.cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stdin pipe in startWebcamRecording: %v", err)
	}

	err = wr.cmd.Start()
	if err != nil {
		return nil, err
	}

	//* looks like this is not accurate, so using file time again instead
	// err = filesystem.WriteCreationTimeUsingNow(filePath)
	// if err != nil {
	// 	fmt.Printf("failed to write creation timestamp for %s: %v\n", rtspPath, err)
	// }

	wr.log("capture started")

	ccKey := config.CC_TOP_CAM
	if rtspPath == *config.FrontCamRtspPath {
		ccKey = config.CC_FRONT_CAM
	}
	SaveCropConfig(ccKey, filePath)

	return wr, nil
}

func (wr *webcamRecorder) Stop() {
	// ensure footage stops after "paused" state report is saved
	time.Sleep(time.Millisecond * time.Duration(500))

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

	err = filesystem.WriteCreationTimeUsingMetadata(wr.filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (wr *webcamRecorder) log(s string, args ...interface{}) {
	prefix := fmt.Sprintf("%s ", wr.name)
	fmt.Printf(prefix+s+"\n", args...)
}
