package obs

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/events"
	"github.com/gkstretton/asol-protos/go/topics_backend"
	"github.com/gkstretton/dark/services/goo/mqtt"
	"github.com/gkstretton/dark/services/goo/session"
	"github.com/gorilla/websocket"
)

var c *goobs.Client
var lock sync.Mutex

const retryWaitS = 5

var sm *session.SessionManager

func Run(s *session.SessionManager) {
	sm = s
	fmt.Println("Running OBS controller")

	go sessionListener(s)
	go connectionListener(s)

	mqtt.Subscribe(topics_backend.TOPIC_STREAM_START, startStream)
	mqtt.Subscribe(topics_backend.TOPIC_STREAM_END, endStream)

	mqtt.Subscribe(topics_backend.TOPIC_STREAM_STATUS_GET, func(topic string, payload []byte) {
		publishStreamStatus()
	})

	publishStreamStatus()
}

func connectionListener(sm *session.SessionManager) {
	reconnect := make(chan bool)
	for {
		fmt.Printf("Attempting connection to OBS @ %s...\n", os.Getenv("OBS_LANDSCAPE_URL"))

		lock.Lock()
		c = nil
		lock.Unlock()

		var err error
		var newClient *goobs.Client
		newClient, err = goobs.New(os.Getenv("OBS_LANDSCAPE_URL"))
		for err != nil {
			fmt.Printf("failed to create obs ws client, retrying in %d seconds: %v\n", retryWaitS, err)
			time.Sleep(time.Second * time.Duration(retryWaitS))
			newClient, err = goobs.New(os.Getenv("OBS_LANDSCAPE_URL"))
		}

		lock.Lock()
		c = newClient
		lock.Unlock()

		resp, err := c.General.GetVersion()
		if err != nil {
			fmt.Printf("failed to get obs version on connect: %v\n", err)
			continue
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
					publishStreamStatus()
					reconnect <- true
				} else {
					fmt.Printf("misc obs error: %v\n", innerErr)
					publishStreamStatus()
				}
			}
			state, ok := i.(*events.StreamStateChanged)
			if ok {
				fmt.Printf("stream state changed to %t\n", state.OutputActive)
				publishStreamStatus()
			}
		})

		handleSessionEvent(sm)
		setCropConfig()

		<-reconnect
	}
}
