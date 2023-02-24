package obs

import (
	"fmt"
	"testing"
	"time"

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
	time.Sleep(time.Second * 5)

	err := setScene("error")
	assert.NoError(t, err)
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
