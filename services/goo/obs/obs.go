package obs

import (
	"fmt"
	"os"
	"time"

	"github.com/andreykaipov/goobs"
	"github.com/gkstretton/dark/services/goo/config"
	"github.com/gkstretton/dark/services/goo/mqtt"
	"github.com/gkstretton/dark/services/goo/session"
	"github.com/gorilla/websocket"
)

var c *goobs.Client

const retryWaitS = 5

func Run(sm *session.SessionManager) {
	fmt.Println("Running OBS controller")

	go connectionListener()
	go sessionListener(sm)

	mqtt.Subscribe(config.TOPIC_STREAM_START, startStream)
	mqtt.Subscribe(config.TOPIC_STREAM_END, endStream)
}

func connectionListener() {
	var err error
	c, err = goobs.New(os.Getenv("OBS_LANDSCAPE_URL"))
	for err != nil {
		fmt.Printf("failed to create obs ws client, retrying in %d seconds: %v\n", retryWaitS, err)
		time.Sleep(time.Second * time.Duration(retryWaitS))
		c, err = goobs.New(os.Getenv("OBS_LANDSCAPE_URL"))
	}

	c.Conn.SetCloseHandler(func(code int, text string) error {
		fmt.Printf("obs connection closed - %d: %s\n", code, text)

		message := websocket.FormatCloseMessage(code, "")
		c.Conn.WriteControl(websocket.CloseMessage, message, time.Now().Add(time.Second))
		go connectionListener()
		return nil
	})

	resp, err := c.General.GetVersion()
	if err != nil {
		fmt.Printf("failed to get obs version on connect: %v\n", err)
		return
	}
	fmt.Printf("Connected to OBS\n"+
		"\tOBSversion: %s\n"+
		"\tOBSws version: %s\n",
		resp.ObsStudioVersion, resp.ObsWebsocketVersion,
	)
}
