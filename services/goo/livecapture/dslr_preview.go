package livecapture

import (
	"fmt"
	"path/filepath"

	"github.com/gkstretton/dark/services/goo/filesystem"
	"github.com/gkstretton/dark/services/goo/mqtt"
)

const (
	TOPIC_TRIGGER_DSLR      = "asol/dslr-crop-capture"
	TOPIC_TRIGGER_DSLR_RESP = "asol/dslr-crop-capture-resp"
	DSLR_CAPTURE_NAME       = "dslr-preview.jpg"
)

func registerDslrPreviewHandler() {
	mqtt.Subscribe(TOPIC_TRIGGER_DSLR, func(topic string, payload []byte) {
		p := filepath.Join(filesystem.GetBasePath(), DSLR_CAPTURE_NAME)
		err := captureImage(p)
		if err != nil {
			fmt.Printf("error in TRIGGER_DSLR handler: %v\n", err)
		} else {
			err := mqtt.Publish(TOPIC_TRIGGER_DSLR_RESP, "ack")
			if err != nil {
				fmt.Printf("error in Trigger_DSLR handler (resp): %v\n", err)
			}
		}
	})
}
