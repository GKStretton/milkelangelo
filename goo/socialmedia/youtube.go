// Sample Go code for user authorization

package socialmedia

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/dark/services/goo/filesystem"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

const missingClientSecretsMessage = `
Please configure OAuth 2.0
`

// getYoutubeClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getYoutubeClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := youtubeTokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
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

// youtubeTokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func youtubeTokenCacheFile() (string, error) {
	return filepath.Join(filesystem.GetBasePath(), "kv", "youtube-credentials.json"), nil
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

const youtubeChannelId = "UCAvwN8vS3f0FFNxWoh3jqyQ"

// can't get this to work
// 2023/11/18 09:56:06 Error making API call: googleapi: Error 400: Each comment thread must be linked to a channel or video.<ul><li>If the comment applies to a channel, make sure that the resource specified in the request body provides a value for the <code><a href="/youtube/v3/docs/commentThreads#snippet.channelId">snippet.channelId</a></code> property. A comment that applies to a channel appears on the channel's <b>Discussion</b> tab.</li><li>If the comment applies to a video, make sure the resource specifies values for both the <code><a href="/youtube/v3/docs/commentThreads#snippet.channelId">snippet.channelId</a></code> and <code><a href="/youtube/v3/docs/commentThreads#snippet.videoId">snippet.videoId</a></code> properties. A comment that applies to a video appears on the video's watch page.</li></ul>, channelOrVideoIdMissing
// note: channel comments (discussion  tab) is deprecated
// then there was the activities / channel bulletins, but that's deprecated too
// ! so I don't think we can do posts with the api...
func doPost(service *youtube.Service) {
	call := service.CommentThreads.Insert([]string{"snippet"}, &youtube.CommentThread{
		Snippet: &youtube.CommentThreadSnippet{
			// CanReply: true,
			ChannelId: youtubeChannelId,
			// IsPublic: true,
			// VideoId:   "7UvSoNPQd3A",
			TopLevelComment: &youtube.Comment{
				Snippet: &youtube.CommentSnippet{
					// ChannelId:    channelId,
					TextOriginal: "api test comment",
				},
			},
		},
	})
	_, err := call.Do()
	handleYoutubeError(err, "")
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
	client := getYoutubeClient(ctx, config)
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

func (m *youtubeManager) Upload(req *UploadRequest) (string, error) {
	if req.Type != VideoPost {
		return "", fmt.Errorf("youtube api doesn't support type %d", req.Type)
	}
	upload := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title:       req.Title,
			Description: req.Description,
		},
	}
	if req.Unlisted {
		upload.Status = &youtube.VideoStatus{
			PrivacyStatus: "unlisted",
		}
	}

	call := m.s.Videos.Insert([]string{"snippet", "status"}, upload)

	file, err := os.Open(req.ContentFilePath)
	defer file.Close()
	if err != nil {
		log.Fatalf("Error opening %v: %v", file, err)
	}
	resp, err := call.Media(file).Do()
	handleYoutubeError(err, "")
	fmt.Printf("Youtube upload successful! Video ID: %v\n", resp.Id)
	return "https://youtube.com/watch?v=" + resp.Id, nil
}

func TestYoutubeClient() {
	m := newYoutubeManager()
	s, err := m.Upload(&UploadRequest{
		Platform:          machinepb.SocialPlatform_SOCIAL_PLATFORM_YOUTUBE,
		Type:              VideoPost,
		Title:             "api test title",
		Description:       "api test desc",
		ContentFilePath:   "/mnt/md0/light-stores/session_content/latest_production/video/post/CONTENT_TYPE_SHORTFORM-overlay.0.mp4",
		ThumbnailFilePath: "",
	})
	fmt.Println(err)
	fmt.Println(s)
}
