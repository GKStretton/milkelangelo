package socialmedia

import (
	"fmt"

	"github.com/gkstretton/asol-protos/go/machinepb"
)

type ContentType int

const (
	VideoPost ContentType = iota
	ImagePost
	TextOnlyPost
)

type UploadRequest struct {
	Platform          machinepb.SocialPlatform
	Type              ContentType
	Title             string
	Description       string
	ContentFilePath   string
	ThumbnailFilePath string
}

type platformManager interface {
	// returns url of upload
	Upload(*UploadRequest) (string, error)
}

type SocialManager struct {
	platformManagers map[machinepb.SocialPlatform]platformManager
}

func NewSocialManager() *SocialManager {
	return &SocialManager{
		platformManagers: map[machinepb.SocialPlatform]platformManager{
			machinepb.SocialPlatform_SOCIAL_PLATFORM_INSTAGRAM: newInstagramManager(),
		},
	}
}

func (s *SocialManager) Upload(req *UploadRequest) (url string, err error) {
	pm, ok := s.platformManagers[req.Platform]
	if !ok || pm == nil {
		return "", fmt.Errorf("no platform manager for %s", req.Platform)
	}

	fmt.Printf("Upload attempt: %s %s\n", req.ContentFilePath, req.ThumbnailFilePath)
	return pm.Upload(req)
}
