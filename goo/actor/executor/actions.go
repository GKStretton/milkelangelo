package executor

import (
	"fmt"
	"sync"

	"github.com/gkstretton/asol-protos/go/topics_firmware"
	"github.com/gkstretton/dark/services/goo/mqtt"
)

// goTo functionality is treated as a resource. For example if something is
// running goto and dispense, a voting round should respect that and not send
// goTo commands.
var goToLock sync.Mutex

func collect(vialNo int, volUl int) {
	mqtt.Publish(topics_firmware.TOPIC_COLLECT, fmt.Sprintf("%d,%d", vialNo, volUl))
}

func goTo(x, y float32) {
	mqtt.Publish(topics_firmware.TOPIC_GOTO_XY, fmt.Sprintf("%.3f,%.3f", x, y))
}

func dispense() error {
	return mqtt.Publish(
		topics_firmware.TOPIC_DISPENSE,
		fmt.Sprintf("%.1f", getDispenseVolume()),
	)
}
