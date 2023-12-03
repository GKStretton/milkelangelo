package executor

import (
	"github.com/gkstretton/dark/services/goo/events"
	"github.com/gkstretton/dark/services/goo/vialprofiles"
)

func getDispenseVolume() float32 {
	sr := events.GetLatestStateReportCopy()

	const fallbackVolume float32 = 15
	profile := vialprofiles.GetSystemVialProfile(int(sr.PipetteState.VialHeld))

	if profile == nil {
		l.Printf("error getting dispense volume, using fallback %.1f\n", fallbackVolume)
		return fallbackVolume
	}

	return profile.DispenseVolumeUl
}
