package contentscheduler

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gkstretton/asol-protos/go/machinepb"
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

	for contentType, contentTypeStatus := range plan.ContentStatuses {
		for _, post := range contentTypeStatus.Posts {
			ct := machinepb.ContentType(machinepb.ContentType_value[contentType])
			if !post.Uploaded && postIsDue(post) {
				url, err := m.handlePost(ct, post, sessionNumber)
				if err != nil {
					fmt.Printf("failed to upload %s to %s: %v\n", contentType, post.Platform, err)
					continue
				}
				post.Url = url
				post.Uploaded = true
			}
		}
	}

	return writeContentPlanToFile(path, plan)
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
