package filesystem

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gkstretton/dark/services/goo/config"
)

var (
	basePath                    = flag.String("basePath", config.LightStores(), "base path for storage data")
	contentPath                 = flag.String("sessionPath", "session_content", "path for session content")
	metadataPath                = flag.String("metadataPath", "session_metadata", "path for session metadata")
	rawVideoPath                = flag.String("rawVideoPath", "video/raw", "path within session, of raw video")
	stateReportFileName         = flag.String("stateReportFileName", "state-reports.yml", "filename for list of state reports")
	contentPlanFileName         = flag.String("contentPlanFileName", "content_plan.yml", "filename for content plan")
	dispenseMetadataFileName    = flag.String("dispenseMetadataFileName", "dispense-metadata.yml", "filename for dispense metadata")
	vialProfileSnapshotFileName = flag.String("vialProfileSnapshotFileName", "vial-profiles.yml", "filename for vial profiles")
)

func AssertBasePaths() {
	e := Exists(*basePath)
	if !e {
		panic("base path '" + *basePath + "' does not exist")
	}
}

func GetBasePath() string {
	return *basePath
}

func Exists(filepath string) bool {
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
	err := os.MkdirAll(p, 0777)
	if err != nil {
		panic(fmt.Errorf("failed to create metadata path: %v", err))
	}
	SetPerms(p)
	return p
}

func GetStateReportPath(sessionId uint64) string {
	p := filepath.Join(
		*basePath,
		*contentPath,
		strconv.Itoa(int(sessionId)),
	)
	err := os.MkdirAll(p, 0777)
	if err != nil {
		panic(fmt.Errorf("failed to create state report path: %v", err))
	}
	SetPerms(p)
	return filepath.Join(p, *stateReportFileName)
}

func GetContentPlanPath(sessionId uint64) string {
	p := filepath.Join(
		*basePath,
		*contentPath,
		strconv.Itoa(int(sessionId)),
	)
	err := os.MkdirAll(p, 0777)
	if err != nil {
		panic(fmt.Errorf("failed to create state report path: %v", err))
	}
	SetPerms(p)
	return filepath.Join(p, *contentPlanFileName)
}

func GetDispenseMetadataPath(sessionId uint64) string {
	p := filepath.Join(
		*basePath,
		*contentPath,
		strconv.Itoa(int(sessionId)),
	)
	err := os.MkdirAll(p, 0777)
	if err != nil {
		panic(fmt.Errorf("failed to create state report path: %v", err))
	}
	SetPerms(p)
	return filepath.Join(p, *dispenseMetadataFileName)
}

func GetVialProfileSnapshotPath(sessionId uint64) string {
	p := filepath.Join(
		*basePath,
		*contentPath,
		strconv.Itoa(int(sessionId)),
	)
	err := os.MkdirAll(p, 0777)
	if err != nil {
		panic(fmt.Errorf("failed to create state report path: %v", err))
	}
	SetPerms(p)
	return filepath.Join(p, *vialProfileSnapshotFileName)
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
	err := os.MkdirAll(p, 0777)
	if err != nil {
		panic(fmt.Errorf("failed to create raw video path: %v", err))
	}
	SetPerms(p)
	return p
}

// GetStillDir mkdirAlls and returns the dir holding the still frame
// selected from the session's top-cam recording (selected.jpg, selected.yml).
func GetStillDir(sessionId uint64) string {
	p := filepath.Join(
		*basePath,
		*contentPath,
		strconv.Itoa(int(sessionId)),
		"still",
	)
	err := os.MkdirAll(p, 0777)
	if err != nil {
		panic(fmt.Errorf("failed to create still path: %v", err))
	}
	SetPerms(p)
	return p
}

// get dir with post processed videos in
func GetPostVideosDir(sessionId uint64) string {
	return filepath.Join(
		*basePath,
		*contentPath,
		strconv.Itoa(int(sessionId)),
		"video",
		"post",
	)
}

// GetPostStillFile returns the filename of the first .jpg file in the
// generated still dir
func GetPostStillFile(sessionId uint64) (string, error) {
	searchDir := filepath.Join(
		*basePath,
		*contentPath,
		strconv.Itoa(int(sessionId)),
		"still",
		"post",
	)

	files, err := os.ReadDir(searchDir)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".jpg") {
			return filepath.Join(searchDir, file.Name()), nil
		}
	}

	return "", fmt.Errorf("no .jpg files found in directory %s", searchDir)
}

// GetIncrementalFile considers 'outDir' and returns the **full path to** the
// NEXT incremental file name on disk (w/ .'ext'). E.g:
//
//	1.mp4 2.mp4 3.mp4 -> [outDir]/4.mp4
func GetIncrementalFileName(outDir string, ext string) string {
	i := 1
	for {
		p := filepath.Join(outDir, fmt.Sprintf("%04d", i)+"."+ext)
		if !Exists(p) {
			return p
		}
		i++
		if i > 10000 {
			panic("bug in GetIncrementalFileName: filename should likely not exceed 10000")
		}
	}
}

func GetKeyValueStorePath() string {
	return filepath.Join(*basePath, "kv")
}
