package obs

import (
	"errors"
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
	var reconnect chan bool
	for {
		fmt.Printf("Attempting connection to OBS...\n")
		c = nil
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
		go c.Listen(func(i interface{}) {
			err, ok := i.(error)
			if ok {
				innerErr := errors.Unwrap(err)
				wsErr, ok := innerErr.(*websocket.CloseError)
				if ok {
					fmt.Printf("websocket closed: %v\n", wsErr)
					reconnect <- true
				} else {
					fmt.Printf("misc obs error: %v\n", innerErr)
				}
			}
		})

		handleSessionEvent(sm)
		setCropConfig()

		<-reconnect
	}
}
