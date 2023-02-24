package obs

import (
	"fmt"
	"os"
	"time"

	"github.com/andreykaipov/goobs"
	"github.com/gkstretton/dark/services/goo/config"
	"github.com/gkstretton/dark/services/goo/mqtt"
	"github.com/gkstretton/dark/services/goo/session"
)

var c *goobs.Client

const retryWaitS = 5

var sm *session.SessionManager

func Run(s *session.SessionManager) {
	sm = s
	fmt.Println("Running OBS controller")

	go sessionListener(s)
	go connectionListener(s)

	mqtt.Subscribe(config.TOPIC_STREAM_START, startStream)
	mqtt.Subscribe(config.TOPIC_STREAM_END, endStream)
}

func connectionListener(sm *session.SessionManager) {
	var err error
	c, err = goobs.New(os.Getenv("OBS_LANDSCAPE_URL"))
	for err != nil {
		fmt.Printf("failed to create obs ws client, retrying in %d seconds: %v\n", retryWaitS, err)
		time.Sleep(time.Second * time.Duration(retryWaitS))
		c, err = goobs.New(os.Getenv("OBS_LANDSCAPE_URL"))
	}

	resp, err := c.General.GetVersion()
	if err != nil {
		fmt.Printf("failed to get obs version on connect: %v\n", err)
		return
	}
	fmt.Printf("Connected to OBS\n"+
		"\tOBSversion: %s\n"+
		"\tOBSws version: %s\n",
		resp.ObsVersion, resp.ObsWebSocketVersion,
	)
	handleSessionEvent(sm)
	setCropConfig()
}
