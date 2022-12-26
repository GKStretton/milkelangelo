package filesystem

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

var (
	basePath     = flag.String("basePath", "/mnt/md0/light-stores/", "base path for storage data")
	contentPath  = flag.String("sessionPath", "session_content", "path for session content")
	metadataPath = flag.String("metadataPath", "session_metadata", "path for session metadata")
	rawVideoPath = flag.String("rawVideoPath", "video/raw", "path within session, of raw video")
)

func AssertBasePaths() {
	e := exists(*basePath)
	if !e {
		panic("base path '" + *basePath + "' does not exist")
	}
}

func exists(filepath string) bool {
	_, err := os.Stat(filepath)
	if err == nil {
		return true
	} else if errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		fmt.Printf("unkown Stat error for session base path assertion: %v\n", err)
		return false
	}
}

func GetMetadataDir() string {
	p := filepath.Join(*basePath, *metadataPath)
	err := os.MkdirAll(p, 0744)
	if err != nil {
		panic(fmt.Errorf("failed to create base path: %v", err))
	}
	return p
}

// GetRawVideoDir mkdirAlls the path if it doesn't exist.
//
//	e.g. 5, top-cam
func GetRawVideoDir(sessionId uint64, rtspPath string) string {
	p := filepath.Join(
		*basePath,
		*contentPath,
		strconv.Itoa(int(sessionId)),
		*rawVideoPath,
		rtspPath,
	)
	err := os.MkdirAll(p, 0744)
	if err != nil {
		panic(fmt.Errorf("failed to create base path: %v", err))
	}
	return p
}

// GetIncrementalFile considers 'outDir' and returns the **full path to** the
// next incremental file name on disk (w/ .'ext'). E.g:
//
//	1.mp4 2.mp4 3.mp4 -> [outDir]/4.mp4
func GetIncrementalFileName(outDir string, ext string) string {
	i := 1
	for {
		p := filepath.Join(outDir, strconv.Itoa(i)+"."+ext)
		if !exists(p) {
			return p
		}
		i++
		if i > 10000 {
			panic("bug in GetIncrementalFileName: filename should not likely exceed 10000")
		}
	}
}
