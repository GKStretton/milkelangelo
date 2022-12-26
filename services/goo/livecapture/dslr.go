package livecapture

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/gkstretton/dark/services/goo/filesystem"
)

func setDslrState(b bool) {
	state := "off"
	if b {
		state = "on"
	}
	cmd := exec.Command("./scripts/set-dslr-state.sh", state)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("failed to run set-dslr-state: %v\n", err)
	}
}

func captureImage(sessionId uint64) {
	p := filesystem.GetIncrementalFileName(filesystem.GetRawDslrDir(sessionId), "jpg")

	cmd := exec.Command("./scripts/capture-dslr.sh", p)
	// cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("failed to run capture-dslr: %v\n", err)
	}
}
