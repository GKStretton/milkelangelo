package socialmedia

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gkstretton/dark/services/goo/filesystem"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func RefreshYoutubeCreds() {
	newYoutubeManager()
}

func newYoutubeManager() platformManager {
	ctx := context.Background()

	b, err := os.ReadFile(filepath.Join(filesystem.GetBasePath(), "kv", "youtube_client_secret.json"))
	if err != nil {
		fmt.Printf("Unable to read client secret file: %v\n", err)
		return nil
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/youtube-go-quickstart.json
	config, err := google.ConfigFromJSON(b, youtube.YoutubeUploadScope, youtube.YoutubeForceSslScope)
	if err != nil {
		fmt.Printf("Unable to parse client secret file to config: %v\n", err)
		return nil
	}
	client := getYoutubeHttpClient(ctx, config)
	service, err := youtube.NewService(ctx, option.WithHTTPClient(client))

	err = handleYoutubeError(err, "Error creating YouTube client")
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &youtubeManager{
		s: service,
	}
}

type youtubeManager struct {
	s *youtube.Service
}

// getYoutubeHttpClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getYoutubeHttpClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile := filepath.Join(filesystem.GetBasePath(), "kv", "youtube-credentials-cache.json")

	tok, err := youtubeTokenFromFile(cacheFile)
	if err != nil {
		tok = getYoutubeTokenFromWeb(config)
		saveYoutubeToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

// getYoutubeTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getYoutubeTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// youtubeTokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func youtubeTokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveYoutubeToken uses a file path to create a file and store the
// token in it.
func saveYoutubeToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func handleYoutubeError(err error, message string) error {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		return fmt.Errorf(message+": %v", err.Error())
	}
	return nil
}
