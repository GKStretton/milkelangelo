package socialmedia

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	// This api is great!!!
	// https://github.com/Davincible/goinsta/wiki/01.-Getting-Started
	"github.com/Davincible/goinsta/v3"
	"github.com/gkstretton/dark/services/goo/filesystem"
	"golang.org/x/term"
)

func goinstaConfigPath() string {
	return filepath.Join(filesystem.GetKeyValueStorePath(), ".goinsta")
}

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
	path := goinstaConfigPath()
	insta, err := goinsta.Import(path)
	if err != nil {
		fmt.Printf("failed to load goinsta config '%s': %v\n", path, err)
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

const defaultInstagramUsername = "astudyoflight_"

// RefreshInstagramCreds logs in interactively (username, password, and 2FA code
// if required) and saves the session to the kv store for newInstagramManager.
func RefreshInstagramCreds() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Instagram username [%s]: ", defaultInstagramUsername)
	username, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("failed to read username: %v\n", err)
		return
	}
	username = strings.TrimSpace(username)
	if username == "" {
		username = defaultInstagramUsername
	}

	fmt.Print("Instagram password: ")
	pw, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		fmt.Printf("failed to read password: %v\n", err)
		return
	}

	insta := goinsta.New(username, string(pw))
	err = insta.Login()
	if errors.Is(err, goinsta.Err2FARequired) {
		fmt.Print("2FA code: ")
		code, rerr := reader.ReadString('\n')
		if rerr != nil {
			fmt.Printf("failed to read 2FA code: %v\n", rerr)
			return
		}
		err = insta.TwoFactorInfo.Login2FA(strings.TrimSpace(code))
	}
	if errors.Is(err, goinsta.ErrChallengeRequired) || errors.Is(err, goinsta.ErrCheckpointRequired) {
		fmt.Println("Instagram wants verification: approve the login in the Instagram app (\"This was me\"), then run -insta again.")
	}
	if err != nil {
		fmt.Printf("instagram login failed: %v\n", err)
		return
	}

	path := goinstaConfigPath()
	if err := insta.Export(path); err != nil {
		fmt.Printf("failed to save instagram session to %s: %v\n", path, err)
		return
	}
	fmt.Printf("Saved instagram session to %s\n", path)
}
