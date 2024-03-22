package contentscheduler

import (
	"fmt"
	"path/filepath"

	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/dark/services/goo/filesystem"
	"github.com/gkstretton/dark/services/goo/socialmedia"
)

// If this doesn't return an error, the post will be set to "uploaded = true"
func (m *manager) handlePost(ct machinepb.ContentType, post *machinepb.Post, sessionNumber uint64) (url string, err error) {
	contentPath, thumbnailPath, err := getContentFilePaths(ct, sessionNumber)
	if err != nil {
		return "", fmt.Errorf("could not get post content paths: %v", err)
	}

	if !filesystem.Exists(contentPath + ".completed") {
		return "", fmt.Errorf("content .completed file not found for %s", contentPath)
	}

	req := &socialmedia.UploadRequest{
		Platform:          post.Platform,
		Type:              socialmedia.VideoPost,
		Title:             post.Title,
		Description:       post.Description,
		ContentFilePath:   contentPath,
		ThumbnailFilePath: thumbnailPath,
		Unlisted:          post.Unlisted,
	}
	if ct == machinepb.ContentType_CONTENT_TYPE_STILL {
		req.Type = socialmedia.ImagePost
	}

	return m.s.Upload(req)
}

func getContentFilePaths(ct machinepb.ContentType, sessionNumber uint64) (string, string, error) {
	switch ct {
	case machinepb.ContentType_CONTENT_TYPE_CLEANING,
		machinepb.ContentType_CONTENT_TYPE_SHORTFORM,
		machinepb.ContentType_CONTENT_TYPE_LONGFORM:

		p := filesystem.GetPostVideosDir(sessionNumber)
		basePath := filepath.Join(p, ct.String())
		contentPath := findLatest(basePath, "mp4")
		thumbnailPath := findLatest(basePath+"-thumbnail", "jpg")
		if contentPath == "" {
			return "", "", fmt.Errorf("couldn't get contentPath for %d %s", sessionNumber, ct)
		}

		return contentPath, thumbnailPath, nil
	case machinepb.ContentType_CONTENT_TYPE_DSLR:
		p := filesystem.GetPostVideosDir(sessionNumber)
		basePath := filepath.Join(p, ct.String())
		contentPath := findLatest(basePath, "mp4")
		if contentPath == "" {
			return "", "", fmt.Errorf("couldn't get contentPath for %d %s", sessionNumber, ct)
		}

		return contentPath, "", nil
	case machinepb.ContentType_CONTENT_TYPE_STILL:
		p, err := filesystem.GetPostStillFile(sessionNumber)
		if err != nil {
			return "", "", err
		}
		return p, p, nil
	default:
		return "", "", fmt.Errorf("content type %s is not supported for posting", ct)
	}
}

// returns latest instance of a file like name.n.ext
func findLatest(file, ext string) string {
	latest := ""
	for i := 0; i < 100; i++ {
		p := file + fmt.Sprintf(".%d.%s", i, ext)
		if filesystem.Exists(p) {
			latest = p
		} else {
			return latest
		}
	}
	return latest
}
