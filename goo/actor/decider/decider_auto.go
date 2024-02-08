package decider

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/dark/services/goo/actor/executor"
	"github.com/gkstretton/dark/services/goo/types"
	"github.com/gkstretton/dark/services/goo/util"
	"github.com/gkstretton/dark/services/goo/vialprofiles"
)

type autoDecider struct {
	endTime time.Time
}

func NewAutoDecider(timeout time.Duration) Decider {
	t := time.Now().Add(timeout)
	return &autoDecider{
		endTime: t,
	}
}

// GetRandomVialPos returns a vial position that meets criteria
func GetRandomVialPos() uint64 {
	options := []uint64{}
	snapshot := vialprofiles.GetSystemVialConfigurationSnapshot()
	for i, p := range snapshot.Profiles {
		if p.VialFluid == machinepb.VialProfile_VIAL_FLUID_DYE_WATER_BASED ||
			p.VialFluid == machinepb.VialProfile_VIAL_FLUID_EMULSIFIER {
			options = append(options, i)
		}
	}

	if len(options) == 0 {
		fmt.Println("ERROR: no system vial profiles")
		return 0
	}
	choiceIndex := rand.Intn(len(options))
	return options[choiceIndex]
}

func (d *autoDecider) decideCollection(predictedState *machinepb.StateReport) *types.CollectionDecision {
	return &types.CollectionDecision{
		VialNo:  int(GetRandomVialPos()),
		DropsNo: 3,
	}
}

// decideDispense decides a random location from the unit circle
func (d *autoDecider) decideDispense(predictedState *machinepb.StateReport) *types.DispenseDecision {
	x, y := util.SampleRandomUnitCircleCoordinate()
	return &types.DispenseDecision{
		X: float32(x),
		Y: float32(y),
	}
}

func (d *autoDecider) DecideNextAction(predictedState *machinepb.StateReport) (executor.Executor, error) {
	if predictedState.Status == machinepb.Status_SLEEPING {
		l.Println("invalid state for actor, decided nil.")
		return nil, fmt.Errorf("invalid machine status for decision: %s", predictedState.Status)
	}
	if predictedState.PipetteState.Spent {
		// only end after the dispense is done
		if time.Now().After(d.endTime) {
			l.Println("endTime reached on auto decider, deciding nil.")
			return nil, nil
		}

		l.Println("collection is next, launching decider...")
		decision := d.decideCollection(predictedState)
		return executor.NewCollectionExecutor(decision), nil
	}
	l.Println("dispense is next, launching decider...")
	decision := d.decideDispense(predictedState)
	return executor.NewDispenseExecutor(decision), nil
}
