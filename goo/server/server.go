// package server is an http server
package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gkstretton/dark/services/goo/config"
	"github.com/gkstretton/dark/services/goo/filesystem"
	"gopkg.in/yaml.v3"
)

func httpErr(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}

func cors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func getSessionId(r *http.Request) (uint64, error) {
	sessionId := r.URL.Query().Get("session_id")
	if sessionId == "" {
		return 0, errors.New("session_id query parameter not specified")
	}
	sessionIdInt, err := strconv.Atoi(sessionId)
	if err != nil {
		return 0, err
	}
	return uint64(sessionIdInt), nil
}

// getVideoPath returns the path of the named top-cam recording in this
// session, rejecting anything that isn't a plain file name.
func getVideoPath(sessionId uint64, r *http.Request) (string, error) {
	name := r.URL.Query().Get("name")
	if name == "" {
		return "", errors.New("name query parameter not specified")
	}
	if name != filepath.Base(name) || !strings.HasSuffix(name, ".mp4") {
		return "", fmt.Errorf("invalid video name '%s'", name)
	}
	return filepath.Join(filesystem.GetRawVideoDir(sessionId, *config.TopCamRtspPath), name), nil
}

func Start() {
	go func() {
		http.HandleFunc("/list-session-videos", func(w http.ResponseWriter, r *http.Request) {
			cors(w)
			sessionId, err := getSessionId(r)
			if err != nil {
				httpErr(w, err)
				return
			}

			names, err := listSessionVideos(sessionId)
			if err != nil {
				httpErr(w, err)
				return
			}

			e := json.NewEncoder(w)
			e.Encode(names)
		})
		http.HandleFunc("/get-session-video", func(w http.ResponseWriter, r *http.Request) {
			cors(w)
			sessionId, err := getSessionId(r)
			if err != nil {
				httpErr(w, err)
				return
			}
			p, err := getVideoPath(sessionId, r)
			if err != nil {
				httpErr(w, err)
				return
			}

			// ServeFile supports range requests, needed for seeking
			http.ServeFile(w, r, p)
		})
		http.HandleFunc("/select-still", func(w http.ResponseWriter, r *http.Request) {
			cors(w)
			sessionId, err := getSessionId(r)
			if err != nil {
				httpErr(w, err)
				return
			}
			p, err := getVideoPath(sessionId, r)
			if err != nil {
				httpErr(w, err)
				return
			}
			t, err := strconv.ParseFloat(r.URL.Query().Get("t"), 64)
			if err != nil {
				httpErr(w, fmt.Errorf("invalid t query parameter: %v", err))
				return
			}

			err = selectStill(sessionId, p, t)
			if err != nil {
				httpErr(w, err)
				return
			}
			w.WriteHeader(http.StatusOK)
		})
		http.HandleFunc("/get-still", func(w http.ResponseWriter, r *http.Request) {
			cors(w)
			sessionId, err := getSessionId(r)
			if err != nil {
				httpErr(w, err)
				return
			}

			w.Header().Set("Cache-Control", "no-store")
			http.ServeFile(w, r, filepath.Join(filesystem.GetStillDir(sessionId), "selected.jpg"))
		})
		log.Fatal(http.ListenAndServe(":8089", nil))
	}()
}

// listSessionVideos returns the top-cam recordings for this session, e.g.
// ["0001.mp4", "0002.mp4"] if the session was paused and resumed.
func listSessionVideos(sessionId uint64) ([]string, error) {
	p := filesystem.GetRawVideoDir(sessionId, *config.TopCamRtspPath)
	files, err := os.ReadDir(p)
	if err != nil {
		return nil, err
	}
	names := []string{}
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".mp4") {
			names = append(names, f.Name())
		}
	}

	return names, nil
}

type selectedStill struct {
	Video         string  `yaml:"video"`
	VideoTimeS    float64 `yaml:"video_time_s"`
	TimestampUnix float64 `yaml:"timestamp_unix"`
}

// selectStill extracts the frame at t seconds into videoPath, cropped with the
// recording's crop config, to still/selected.jpg. The absolute timestamp of
// the frame is written to still/selected.yml.
func selectStill(sessionId uint64, videoPath string, t float64) error {
	rawStart, err := os.ReadFile(videoPath + ".creationtime")
	if err != nil {
		return fmt.Errorf("no creationtime for %s, is it still recording? %v", videoPath, err)
	}
	start, err := strconv.ParseFloat(strings.TrimSpace(string(rawStart)), 64)
	if err != nil {
		return fmt.Errorf("failed to parse creationtime for %s: %v", videoPath, err)
	}

	args := []string{"-y", "-ss", strconv.FormatFloat(t, 'f', 3, 64), "-i", videoPath, "-frames:v", "1", "-q:v", "2"}
	rawCrop, err := os.ReadFile(videoPath + ".yml")
	if err == nil {
		cc := config.CropConfig{}
		if err := yaml.Unmarshal(rawCrop, &cc); err != nil {
			return fmt.Errorf("failed to parse crop config for %s: %v", videoPath, err)
		}
		args = append(args, "-vf", fmt.Sprintf("crop=%d:%d:%d:%d",
			cc.RightAbs-cc.LeftAbs, cc.BottomAbs-cc.TopAbs, cc.LeftAbs, cc.TopAbs))
	} else {
		fmt.Printf("no crop config for %s, still will be uncropped: %v\n", videoPath, err)
	}

	dir := filesystem.GetStillDir(sessionId)
	imagePath := filepath.Join(dir, "selected.jpg")
	args = append(args, imagePath)
	out, err := exec.Command("ffmpeg", args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg frame extraction failed: %v, %s", err, string(out))
	}
	filesystem.SetPerms(imagePath)

	data, err := yaml.Marshal(selectedStill{
		Video:         filepath.Base(videoPath),
		VideoTimeS:    t,
		TimestampUnix: start + t,
	})
	if err != nil {
		return err
	}
	ymlPath := filepath.Join(dir, "selected.yml")
	if err := os.WriteFile(ymlPath, data, 0666); err != nil {
		return err
	}
	filesystem.SetPerms(ymlPath)
	return nil
}
