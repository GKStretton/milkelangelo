package decider

import (
	"log"
	"os"

	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/dark/services/goo/actor/executor"
	"github.com/gkstretton/dark/services/goo/vialprofiles"
)

var l = log.New(os.Stdout, "[decider] ", log.Flags())

type Decider interface {
	DecideCollection(predictedState *machinepb.StateReport) executor.Executor
	DecideDispense(predictedState *machinepb.StateReport) executor.Executor
}

func getVialVolume(vialNo int) float32 {
	const fallbackVolume float32 = 15
	profile := vialprofiles.GetSystemVialProfile(vialNo)

	if profile == nil {
		l.Printf("error getting vial volume, using fallback %.1f\n", fallbackVolume)
		return fallbackVolume
	}

	return profile.DispenseVolumeUl
}
