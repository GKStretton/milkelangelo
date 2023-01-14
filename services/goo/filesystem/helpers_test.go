package filesystem

import (
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWriteCreationTime(t *testing.T) {
	os.Chdir("../")
	start := time.Now()

	filePath := "/tmp/test"
	err := os.WriteFile(filePath, []byte("testfile"), 0666)
	assert.NoError(t, err)

	err = WriteCreationTime(filePath)
	assert.NoError(t, err)

	b, err := os.ReadFile(filePath + ".creationtime")
	assert.NoError(t, err)
	s := strings.Split(string(b), ".")
	assert.Len(t, s, 2)
	t.Log(s)

	seconds, err := strconv.Atoi(s[0])
	assert.NoError(t, err)
	nanos, err := strconv.Atoi(s[1])
	assert.NoError(t, err)

	t.Logf("seconds: %d\n", seconds)
	t.Logf("nanos: %d\n", nanos)

	createdTime := time.Unix(int64(seconds), int64(nanos))
	t.Logf("start: %s\n", start.String())
	t.Logf("created: %s\n", createdTime.String())

	msDiff := createdTime.UnixMilli() - start.UnixMilli()
	if msDiff < 0 {
		msDiff *= -1
	}
	assert.True(t, msDiff < 500)

	err = os.Remove(filePath)
	assert.NoError(t, err)
	err = os.Remove(filePath + ".creationtime")
	assert.NoError(t, err)
}
