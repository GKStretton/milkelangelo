package executor

import (
	"fmt"

	"github.com/gkstretton/asol-protos/go/topics_firmware"
	"github.com/gkstretton/dark/services/goo/mqtt"
)

func collect(vialNo int, volUl int) {
	go mqtt.Publish(topics_firmware.TOPIC_COLLECT, fmt.Sprintf("%d,%d", vialNo, volUl))
}

func goTo(x, y float32) {
	go mqtt.Publish(topics_firmware.TOPIC_GOTO_XY, fmt.Sprintf("%.3f,%.3f", x, y))
}

func dispense() error {
	return mqtt.Publish(
		topics_firmware.TOPIC_DISPENSE,
		fmt.Sprintf("%.1f", getDispenseVolume()),
	)
}
