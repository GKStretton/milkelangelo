// Sample Go code for user authorization

package socialmedia

import (
	"fmt"
	"os"

	"github.com/gkstretton/asol-protos/go/machinepb"
	"google.golang.org/api/youtube/v3"
)

func (m *youtubeManager) Upload(req *UploadRequest) (string, error) {
	if req.Type != VideoPost {
		return "", fmt.Errorf("youtube api doesn't support type %d", req.Type)
	}
	upload := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			ChannelId:   youtubeChannelId,
			Title:       req.Title,
			Description: req.Description,
		},
		Status: &youtube.VideoStatus{
			PrivacyStatus:           "public",
			SelfDeclaredMadeForKids: false,
		},
	}
	if req.Unlisted {
		upload.Status.PrivacyStatus = "unlisted"
	}

	call := m.s.Videos.Insert([]string{"snippet", "status"}, upload)

	file, err := os.Open(req.ContentFilePath)
	if err != nil {
		return "", fmt.Errorf("error opening %s: %v", file.Name(), err)
	}
	defer file.Close()
	resp, err := call.Media(file).Do()
	err = handleYoutubeError(err, "")
	if err != nil {
		return "", err
	}
	fmt.Printf("Youtube upload successful! Video ID: %s\n", resp.Id)

	if req.ThumbnailFilePath != "" {
		err = m.setThumbnail(resp.Id, req.ThumbnailFilePath)
		if err != nil {
			fmt.Printf("error setting thumbnail: %v\n", err)
		}
	}

	return "https://youtube.com/watch?v=" + resp.Id, nil
}

func (m *youtubeManager) setThumbnail(videoId string, thumbnailPath string) error {
	tnCall := m.s.Thumbnails.Set(videoId)

	tn, err := os.Open(thumbnailPath)
	if err != nil {
		return fmt.Errorf("error opening thumbnail %s: %v", thumbnailPath, err)
	}
	defer tn.Close()

	_, err = tnCall.Media(tn).Do()
	if err != nil {
		return fmt.Errorf("error setting thumbnail: %v", err)
	}
	fmt.Printf("Thumbnail set for video ID: %s\n", videoId)
	return nil
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
		Unlisted:          true,
	})
	fmt.Println(err)
	fmt.Println(s)
}

//! experimental, not working

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
