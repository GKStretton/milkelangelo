package filesystem

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func CreateSymlink(original, new string) error {
	if err := removeSymlink(new); err != nil {
		return fmt.Errorf("failed to unlink latest: %v", err)
	}
	if err := os.Symlink(original, new); err != nil {
		return fmt.Errorf("failed to symlink latest: %v", err)
	}
	SetPerms(new)
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

func WriteCreationTimeUsingMetadata(filePath string) error {
	cmd := exec.Command("./scripts/get-creation-timestamp.sh", filePath)
	var ts bytes.Buffer
	cmd.Stdout = &ts
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error getting creation time for %s: %v", filePath, err)
	}
	// discard \n
	ts.Truncate(ts.Len() - 1)

	// fmt.Printf("got creation timestamp %s\n", ts.Bytes())
	name := filePath + ".creationtime"
	if err := os.WriteFile(name, ts.Bytes(), 0666); err != nil {
		return fmt.Errorf("error writing creation time for %s: %v", filePath, err)
	}
	SetPerms(name)
	return nil
}

func WriteCreationTimeUsingNow(filePath string) error {
	now := time.Now()
	ts := fmt.Sprintf("%d.%d", now.Unix(), now.Nanosecond())

	name := filePath + ".creationtime"
	if err := os.WriteFile(name, []byte(ts), 0666); err != nil {
		return fmt.Errorf("error writing creation time for %s: %v", filePath, err)
	}
	SetPerms(name)
	return nil
}

func SetPerms(p string) {
	// chown to 1000:1000 (host user)
	// just do everything to be sure
	if err := chownRecursive(GetBasePath()+*contentPath, 1000, 1000); err != nil {
		fmt.Printf("failed to chown %s: %v\n", p, err)
	}
	if err := chownRecursive(p, 1000, 1000); err != nil {
		fmt.Printf("failed to chown %s: %v\n", p, err)
	}
}

func chownRecursive(path string, uid, gid int) error {
	cmd := exec.Command("chown", "-R", fmt.Sprintf("%d:%d", uid, gid), path)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("chown failed: %v, %s", err, string(out))
	}
	return nil
}
