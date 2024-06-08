package executor

import (
	"github.com/gkstretton/dark/services/goo/events"
	"github.com/gkstretton/dark/services/goo/vialprofiles"
)

func getDispenseVolume() float32 {
	sr := events.GetLatestStateReportCopy()

	return getVialDropVolume(int(sr.PipetteState.VialHeld))
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
