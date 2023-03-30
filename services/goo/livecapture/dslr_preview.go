package livecapture

import (
	"fmt"
	"path/filepath"

	"github.com/gkstretton/asol-protos/go/topics_backend"
	"github.com/gkstretton/dark/services/goo/filesystem"
	"github.com/gkstretton/dark/services/goo/mqtt"
)

const (
	DSLR_CAPTURE_NAME = "dslr-preview.jpg"
)

func registerDslrPreviewHandler() {
	mqtt.Subscribe(topics_backend.TOPIC_TRIGGER_DSLR, func(topic string, payload []byte) {
		p := filepath.Join(filesystem.GetBasePath(), DSLR_CAPTURE_NAME)
		pl := "ack"
		err := captureImage(p)
		if err != nil {
			fmt.Printf("error in TRIGGER_DSLR handler: %v\n", err)
			pl = "err"
		}
		err = mqtt.Publish(topics_backend.TOPIC_TRIGGER_DSLR_RESP, pl)
		if err != nil {
			fmt.Printf("error in Trigger_DSLR handler (resp): %v\n", err)
		}
	})
}
