package obs

import (
	"testing"

	"github.com/andreykaipov/goobs"
	"github.com/stretchr/testify/assert"
)

func initTestClient(t *testing.T) {
	cli, err := goobs.New("localhost:4444")
	assert.NoError(t, err)

	c = cli
}

func TestSetSession(t *testing.T) {
	initTestClient(t)

	setSessionNumber(3)
}

func TestSetScene(t *testing.T) {
	initTestClient(t)

	setScene("fallback")
}

func TestSetCropConfig(t *testing.T) {
	initTestClient(t)

	setCropConfig()
}
