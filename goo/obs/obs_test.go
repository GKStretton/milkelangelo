package obs

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/inputs"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func initTestClient(t *testing.T) {
	cli, err := goobs.New("192.168.0.37:4444")
	assert.NoError(t, err)

	c = cli
}

// func TestSetSession(t *testing.T) {
// 	initTestClient(t)

// 	setSessionNumber(3)
// }

func TestEventListener(t *testing.T) {
	initTestClient(t)
	go c.Listen(func(i interface{}) {
		err, ok := i.(error)
		if ok {
			innerErr := errors.Unwrap(err)
			wsErr, ok := innerErr.(*websocket.CloseError)
			if ok {
				t.Logf("websocket closed: %v", wsErr)
			} else {
				t.Logf("misc error: %v", innerErr)
			}
		}

	})
	time.Sleep(time.Second * 10)
	t.Log("test")
}

func TestGetStreamStatus(t *testing.T) {
	initTestClient(t)

	fmt.Println(isStreamLive())
}

func TestSetScene(t *testing.T) {
	initTestClient(t)

	err := setScene("complete")
	assert.NoError(t, err)
}

func TestSetSessionNumber(t *testing.T) {
	initTestClient(t)

	setSessionNumber(3548, true)
}

func TestStartStream(t *testing.T) {
	initTestClient(t)

	startStream("", []byte{})
}

func TestSetCropConfig(t *testing.T) {
	initTestClient(t)

	setCropConfig()
}

func TestGetInputs(t *testing.T) {
	initTestClient(t)

	resp, err := c.Inputs.GetInputList(&inputs.GetInputListParams{})
	fmt.Println(resp.Inputs[0], err)
}
