package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gkstretton/asol-protos/go/topics_backend"
	"github.com/gkstretton/dark/services/dslrcapture/mqtt"
)

func saveCropConfig(ccKey string, contentPath string) error {
	// e.g. 1.mp4.yml
	ymlPath := contentPath + ".yml"

	mqtt.Publish(topics_backend.TOPIC_KV_GET+ccKey, "")
	fmt.Printf("issuing SubscribeBlocking for %s\n", ccKey)
	config, err := mqtt.SubscribeBlocking(topics_backend.TOPIC_KV_GET_RESP+ccKey, time.Second)
	if err != nil {
		return fmt.Errorf("failed to get cropConfig of %s: %v", ccKey, err)
	}
	fmt.Printf("got cropConfig for dslr capture of %s\n", config)

	err = os.WriteFile(ymlPath, config, 0666)
	if err != nil {
		return fmt.Errorf("failed to write cropConfig of %s to '%s': %v", ccKey, ymlPath, err)
	}
	return nil
}
