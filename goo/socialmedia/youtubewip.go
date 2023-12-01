// Sample Go code for user authorization

package socialmedia

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

const missingClientSecretsMessage = `
Please configure OAuth 2.0
`

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
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

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("youtube-go-quickstart.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func handleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		log.Fatalf(message+": %v", err.Error())
	}
}

func main2() {
	ctx := context.Background()

	b, err := os.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/youtube-go-quickstart.json
	config, err := google.ConfigFromJSON(b, youtube.YoutubeUploadScope, youtube.YoutubeForceSslScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(ctx, config)
	service, err := youtube.New(client)

	handleError(err, "Error creating YouTube client")

	// uploadVideo(service)
	doPost(service)
}

const channelId = "UCAvwN8vS3f0FFNxWoh3jqyQ"

// can't get this to work
// 2023/11/18 09:56:06 Error making API call: googleapi: Error 400: Each comment thread must be linked to a channel or video.<ul><li>If the comment applies to a channel, make sure that the resource specified in the request body provides a value for the <code><a href="/youtube/v3/docs/commentThreads#snippet.channelId">snippet.channelId</a></code> property. A comment that applies to a channel appears on the channel's <b>Discussion</b> tab.</li><li>If the comment applies to a video, make sure the resource specifies values for both the <code><a href="/youtube/v3/docs/commentThreads#snippet.channelId">snippet.channelId</a></code> and <code><a href="/youtube/v3/docs/commentThreads#snippet.videoId">snippet.videoId</a></code> properties. A comment that applies to a video appears on the video's watch page.</li></ul>, channelOrVideoIdMissing
func doPost(service *youtube.Service) {
	call := service.CommentThreads.Insert([]string{"snippet"}, &youtube.CommentThread{
		Snippet: &youtube.CommentThreadSnippet{
			ChannelId: channelId,
			IsPublic:  true,
			// VideoId:   "IYSKEj23bP8",
			TopLevelComment: &youtube.Comment{
				Snippet: &youtube.CommentSnippet{
					ChannelId:    channelId,
					TextOriginal: "api test comment",
				},
			},
		},
	})
	_, err := call.Do()
	handleError(err, "")
}

func uploadVideo(service *youtube.Service) {
	upload := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title:       "api test upload",
			Description: "uploaded via Go client library",
		},
	}
	call := service.Videos.Insert([]string{"snippet", "status"}, upload)

	file, err := os.Open("/mnt/md0/light-stores/session_content/latest_production/video/post/CONTENT_TYPE_SHORTFORM-overlay.0.mp4")
	defer file.Close()
	if err != nil {
		log.Fatalf("Error opening %v: %v", file, err)
	}
	_ = file
	// resp, err := call.Media(file).Do()
	resp, err := call.Do()
	handleError(err, "")
	fmt.Printf("Upload successful! Video ID: %v\n", resp.Id)
}
