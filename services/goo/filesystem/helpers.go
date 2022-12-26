package filesystem

import (
	"errors"
	"fmt"
	"os"
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
