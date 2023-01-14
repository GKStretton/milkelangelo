package filesystem

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
)

func CreateSymlink(original, new string) error {
	if err := removeSymlink(new); err != nil {
		return fmt.Errorf("failed to unlink latest: %v", err)
	}
	if err := os.Symlink(original, new); err != nil {
		return fmt.Errorf("failed to symlink latest: %v", err)
	}
	return nil
}

func removeSymlink(symlinkPath string) error {
	if _, err := os.Lstat(symlinkPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("failed to stat symlink: %v", err)
	}
	if err := os.Remove(symlinkPath); err != nil {
		return fmt.Errorf("failed to remove symlink: %v", err)
	}
	return nil
}

func WriteCreationTime(filePath string) error {
	cmd := exec.Command("./scripts/get-creation-timestamp.sh", filePath)
	var ts bytes.Buffer
	cmd.Stdout = &ts
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error getting creation time for %s: %v", filePath, err)
	}
	// discard \n
	ts.Truncate(ts.Len() - 1)

	fmt.Printf("got creation timestamp %s\n", ts.Bytes())
	if err := os.WriteFile(filePath+".creationtime", ts.Bytes(), 0666); err != nil {
		return fmt.Errorf("error writing creation time for %s: %v", filePath, err)
	}
	return nil
}
