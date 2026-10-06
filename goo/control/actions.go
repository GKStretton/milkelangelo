package control

import (
	"fmt"
	"time"

	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/asol-protos/go/topics_firmware"
	"github.com/gkstretton/dark/services/goo/events"
	"github.com/gkstretton/dark/services/goo/mqtt"
	"github.com/gkstretton/dark/services/goo/vialprofiles"
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

// call dispense, and observe transition (-> dispensing -> not dispensing)
func dispenseBlocking() {
	a1 := events.ConditionWaiter(func(sr *machinepb.StateReport) bool {
		return sr.GetStatus() == machinepb.Status_DISPENSING
	})
	time.Sleep(time.Millisecond * 250)

	err := dispense()
	if err != nil {
		l.Printf("dispense error: %v\n", err)
		innerErr := dispense()
		if innerErr != nil {
			l.Printf("dispense error after retry: %v\n", innerErr)
		}
	}

	l.Println("waiting for DISPENSING")
	<-a1

	a2 := events.ConditionWaiter(func(sr *machinepb.StateReport) bool {
		return sr.GetStatus() != machinepb.Status_DISPENSING
	})
	l.Println("waiting for NOT DISPENSING")
	<-a2
}

func getDispenseVolume() float32 {
	sr := events.GetLatestStateReportCopy()

	return getVialDropVolume(int(sr.GetPipetteState().GetVialHeld()))
}

func getVialDropVolume(vialNo int) float32 {
	const fallbackVolume float32 = 15
	profile := vialprofiles.GetSystemVialProfile(vialNo)

	if profile == nil {
		l.Printf("error getting vial %d volume, using fallback %.1f\n", vialNo, fallbackVolume)
		return fallbackVolume
	}

	return profile.DispenseVolumeUl
}
