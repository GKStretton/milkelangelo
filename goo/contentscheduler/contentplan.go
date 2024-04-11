package contentscheduler

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/dark/services/goo/email"
	"github.com/gkstretton/dark/services/goo/util/protoyaml"
)

func postIsDue(post *machinepb.Post) bool {
	ts := time.Unix(int64(post.ScheduledUnixTimetamp), 0)
	return ts.Before(time.Now())
}

func (m *manager) processContentPlan(path string, sessionNumber uint64) error {
	lock.Lock()
	defer lock.Unlock()

	plan, err := readContentPlanFromFile(path)
	if err != nil {
		return fmt.Errorf("failed to get content plan from file: %v", err)
	}

	if postsStillToUpload(plan) == 0 {
		return nil
	}

	for contentType, contentTypeStatus := range plan.ContentStatuses {
		for _, post := range contentTypeStatus.Posts {
			if post.Uploaded || !postIsDue(post) {
				continue
			}

			ct := machinepb.ContentType(machinepb.ContentType_value[contentType])
			url, err := m.handlePost(ct, post, sessionNumber)
			if err != nil {
				fmt.Printf("failed to upload %s to %s: %v\n", contentType, post.Platform, err)
				continue
			}
			post.Url = url
			post.Uploaded = true
			fmt.Printf("new upload:\n\t%s\n\t%s\n\t%s\n", post.Platform, ct, post.Url)
		}
	}

	err = writeContentPlanToFile(path, plan)
	if err != nil {
		return fmt.Errorf("failed to write content plan to file: %v", err)
	}

	remaining := postsStillToUpload(plan)
	if remaining == 0 {
		sendUploadsCompleteEmail(plan, sessionNumber)
		fmt.Printf("All posts uploaded for session %d\n", sessionNumber)
		return nil
	}

	fmt.Printf("partially processed content plan for session %d (%d posts still to upload)\n", sessionNumber, remaining)
	return nil
}
func sendUploadsCompleteEmail(plan *machinepb.ContentTypeStatuses, sessionNumber uint64) {
	body := "All posts have been uploaded.\n\n"
	for contentType, contentTypeStatus := range plan.ContentStatuses {
		for _, post := range contentTypeStatus.Posts {
			if post.Uploaded {
				body += fmt.Sprintf("%s post uploaded to %s: %s\n", contentType, post.Platform, post.Url)
			}
		}
	}

	err := email.SendEmail(&machinepb.Email{
		Subject:   fmt.Sprintf("Uploads complete for session %d", sessionNumber),
		Body:      body,
		Recipient: machinepb.EmailRecipient_EMAIL_RECIPIENT_SOCIAL_NOTIFICATIONS,
	})
	if err != nil {
		fmt.Printf("failed to send email on all posts uploaded: %v\n", err)
	}
}

func postsStillToUpload(plan *machinepb.ContentTypeStatuses) int {
	c := 0
	for _, contentTypeStatus := range plan.ContentStatuses {
		for _, post := range contentTypeStatus.Posts {
			if !post.Uploaded {
				c++
			}
		}
	}
	return c
}

func readContentPlanFromFile(path string) (*machinepb.ContentTypeStatuses, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	plan := &machinepb.ContentTypeStatuses{}
	err = protoyaml.Unmarshal(b, plan)
	if err != nil {
		return nil, err
	}

	return plan, nil
}

func writeContentPlanToFile(path string, plan *machinepb.ContentTypeStatuses) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := protoyaml.Marshal(plan)
	if err != nil {
		return err
	}

	_, err = f.Write(b)
	return err
}
