// package server is an http server
package server

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gkstretton/dark/services/goo/filesystem"
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

func Start() {
	go func() {
		http.HandleFunc("/list-dslr-post", func(w http.ResponseWriter, r *http.Request) {
			cors(w)
			sessionId := r.URL.Query().Get("session_id")
			if sessionId == "" {
				httpErr(w, errors.New("session_id query parameter not specified"))
				return
			}
			sessionIdInt, err := strconv.Atoi(sessionId)
			if err != nil {
				httpErr(w, err)
				return
			}

			names, err := listSessionDSLRImages(sessionIdInt)
			if err != nil {
				httpErr(w, err)
				return
			}

			e := json.NewEncoder(w)
			e.Encode(names)
		})
		http.HandleFunc("/get-dslr-post", func(w http.ResponseWriter, r *http.Request) {
			cors(w)
			sessionId := r.URL.Query().Get("session_id")
			if sessionId == "" {
				httpErr(w, errors.New("session_id query parameter not specified"))
				return
			}
			sessionIdInt, err := strconv.Atoi(sessionId)
			if err != nil {
				httpErr(w, err)
				return
			}

			imageName := r.URL.Query().Get("image_name")
			if imageName == "" {
				httpErr(w, errors.New("image_name query parameter not specified"))
				return
			}

			fileReader, err := getSessionDSLRImage(sessionIdInt, imageName)
			if err != nil {
				httpErr(w, err)
				return
			}

			io.Copy(w, fileReader)
		})
		http.HandleFunc("/get-dslr-preview", func(w http.ResponseWriter, r *http.Request) {
			cors(w)

			fileReader, err := getPreviewDSLRImage()
			if err != nil {
				httpErr(w, err)
				return
			}

			io.Copy(w, fileReader)
		})
		http.HandleFunc("/select-dslr-post", func(w http.ResponseWriter, r *http.Request) {
			cors(w)
			sessionId := r.URL.Query().Get("session_id")
			if sessionId == "" {
				httpErr(w, errors.New("session_id query parameter not specified"))
				return
			}
			sessionIdInt, err := strconv.Atoi(sessionId)
			if err != nil {
				httpErr(w, err)
				return
			}

			imageName := r.URL.Query().Get("image_name")
			if imageName == "" {
				httpErr(w, errors.New("image_name query parameter not specified"))
				return
			}
			err = selectDslrImage(sessionIdInt, imageName)
			if err != nil {
				httpErr(w, err)
				return
			}
			w.WriteHeader(http.StatusOK)
		})
		log.Fatal(http.ListenAndServe(":8089", nil))
	}()
}

// create selected.jpg symlink for this image
func selectDslrImage(sessionNo int, name string) error {
	p := filesystem.GetPostDslrDir(uint64(sessionNo))
	imagePath := filepath.Join(p, name)
	linkPath := filepath.Join(p, "selected.jpg")
	os.Remove(linkPath)
	return os.Symlink(imagePath, linkPath)
}

func listSessionDSLRImages(sessionNo int) ([]string, error) {
	p := filesystem.GetPostDslrDir(uint64(sessionNo))
	files, err := os.ReadDir(p)
	if err != nil {
		return nil, err
	}
	names := []string{}
	for _, f := range files {
		names = append(names, f.Name())
	}

	return names, nil
}

func getSessionDSLRImage(sessionNo int, name string) (io.Reader, error) {
	p := filesystem.GetPostDslrDir(uint64(sessionNo))
	filename := filepath.Join(p, name)
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func getPreviewDSLRImage() (io.Reader, error) {
	filename := filesystem.GetDslrPreviewFile()
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return f, nil
}
