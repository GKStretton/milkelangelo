package socialmedia

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	// This api is great!!!
	// https://github.com/Davincible/goinsta/wiki/01.-Getting-Started
	"github.com/Davincible/goinsta/v3"
)

const path = "/mnt/md0/light-stores/kv/.goinsta"

type instagramManager struct {
	i *goinsta.Instagram
}

func (m *instagramManager) Upload(req *UploadRequest) (string, error) {
	if req.Type != ImagePost && req.Type != VideoPost {
		return "", fmt.Errorf("instagram doesn't support type %d", req.Type)
	}

	instaReq := &goinsta.UploadOptions{
		Caption: req.Title,
	}

	contentFile, err := os.Open(req.ContentFilePath)
	if err != nil {
		return "", err
	}
	instaReq.File = contentFile

	if req.ThumbnailFilePath != "" {
		thumbnailFile, err := os.Open(req.ThumbnailFilePath)
		if err != nil {
			return "", fmt.Errorf("error reading thumbnail %s: %v", req.ThumbnailFilePath, err)
		}
		instaReq.Thumbnail = thumbnailFile
	}

	item, err := m.i.Upload(instaReq)
	if err != nil {
		return "", err
	}

	url, err := getInstagramUrlFromMediaId(item.GetID())
	if err != nil {
		fmt.Printf("failed to get url for ig post: %v\n", err)
		return "", nil
	}

	fmt.Printf("uploaded ig post %s\n", url)

	return url, nil
}

func newInstagramManager() platformManager {
	insta, err := goinsta.Import(path)
	if err != nil {
		fmt.Printf("failed to load goinsta config: %v\n", err)
		return nil
	}
	defer insta.Export(path)

	return &instagramManager{
		i: insta,
	}
}

// getInstagramUrlFromMediaId converts an Instagram media ID to its corresponding URL.
func getInstagramUrlFromMediaId(mediaID string) (string, error) {
	index := strings.Index(mediaID, "_")
	if index == -1 {
		return "", fmt.Errorf("invalid media ID format")
	}

	numericID, err := strconv.ParseInt(mediaID[:index], 10, 64)
	if err != nil {
		return "", err
	}

	alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
	shortenedID := ""

	for numericID > 0 {
		remainder := numericID % 64
		numericID = (numericID - remainder) / 64
		shortenedID = string(alphabet[remainder]) + shortenedID
	}

	return "https://www.instagram.com/p/" + shortenedID + "/", nil
}

// only needed to generate initial config
func instaLogin() {
	insta := goinsta.New("astudyoflight_", "[password]")

	// Only call Login the first time you login. Next time import your config
	if err := insta.Login(); err != nil {
		panic(err)
	}

	// Export your configuration
	// after exporting you can use Import function instead of New function.
	// insta, err := goinsta.Import("~/.goinsta")
	// it's useful when you want use goinsta repeatedly.
	// Export is deffered because every run insta should be exported at the end of the run
	//   as the header cookies change constantly.
	defer insta.Export(".goinsta")
}
