package obs

import (
	"fmt"
	"testing"

	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/inputs"
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

func TestSetScene(t *testing.T) {
	initTestClient(t)

	err := setScene("complete")
	assert.NoError(t, err)
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
