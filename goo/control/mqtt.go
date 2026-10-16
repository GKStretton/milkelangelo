package control

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gkstretton/asol-protos/go/topics_backend"
	"github.com/gkstretton/asol-protos/go/topics_firmware"
	"github.com/gkstretton/dark/services/goo/mqtt"
	"github.com/gkstretton/dark/services/goo/types"
)

// The actor topics are reused for enabling / disabling control.
func subscribeToBrokerTopics() {
	mqtt.Subscribe(topics_backend.TOPIC_ACTOR_START, func(topic string, payload []byte) {
		l.Println("mqtt enable request")
		SetEnabled(true)
	})
	mqtt.Subscribe(topics_backend.TOPIC_ACTOR_STOP, func(topic string, payload []byte) {
		l.Println("mqtt disable request")
		SetEnabled(false)
	})
	mqtt.Subscribe(topics_backend.TOPIC_ACTOR_STATUS_GET, func(topic string, payload []byte) {
		publishStatus()
	})

	// follow bowl speed requests from anywhere (this, the interface) so
	// clients see the latest
	mqtt.Subscribe(topics_firmware.TOPIC_SET_BOWL_STEPS_PER_SEC, func(topic string, payload []byte) {
		speed, err := parseBowlSpeed(string(payload))
		if err != nil {
			l.Printf("bad bowl speed request %q: %v\n", payload, err)
			return
		}
		if stateUpdater != nil {
			stateUpdater.UpdateState(func(state *types.GooState) {
				state.BowlSpeed = speed
			})
		}
	})

	publishStatus()
}

// parseBowlSpeed reads the speed from a "speed,accel" bowl request
func parseBowlSpeed(payload string) (float32, error) {
	speed, _, _ := strings.Cut(payload, ",")
	f, err := strconv.ParseFloat(strings.TrimSpace(speed), 32)
	return float32(f), err
}

func publishStatus() {
	mqtt.Publish(topics_backend.TOPIC_ACTOR_STATUS_RESP, fmt.Sprintf("%t", IsEnabled()))
}
