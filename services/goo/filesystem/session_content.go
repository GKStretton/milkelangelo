package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

// InitSessionContent creates a session folder on disk that content uses
// It also sets the 'latest' symlink to point to this session folder
func InitSessionContent(sessionId uint64) error {
	p := filepath.Join(*basePath, *contentPath, strconv.Itoa(int(sessionId)))
	if err := os.Mkdir(p, 0744); err != nil {
		return fmt.Errorf("failed to mkdir: %v", err)
	}

	latestPath := filepath.Join(*basePath, *contentPath, "latest")
	if err := CreateSymlink(p, latestPath); err != nil {
		return fmt.Errorf("failed to symlink latest session folder: %v", err)
	}

	return nil
}
