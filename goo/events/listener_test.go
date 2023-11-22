package events

import "testing"

func TestAppendFailedDispense(t *testing.T) {
	appendFailedDispense(58, 0, 0)
	appendFailedDispense(58, 0, 1)
}
