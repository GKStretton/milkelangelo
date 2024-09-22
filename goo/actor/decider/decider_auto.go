package decider

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/dark/services/goo/actor/executor"
	"github.com/gkstretton/dark/services/goo/types"
	"github.com/gkstretton/dark/services/goo/vialprofiles"
)

type autoDecider struct {
	endTime        time.Time
	rand           *rand.Rand
	testing        bool
	emulsifierUsed bool
}

func NewAutoDecider(timeout time.Duration, seed int64, testing bool) Decider {
	t := time.Now().Add(timeout)
	seededRand := rand.New(rand.NewSource(seed))
	return &autoDecider{
		endTime: t,
		rand:    seededRand,
		testing: testing,
	}
}

// GetRandomVialPos returns a vial position that meets criteria
func (d *autoDecider) GetRandomVialPos(predictedState *machinepb.StateReport) uint64 {
	if d.testing {
		// pos 2 is empty
		return 2
	}

	// only allow one emul, and ensure it's not too early (>= 1 means only from second collection)
	emulsifierAllowed := predictedState.CollectionRequest.RequestNumber >= 3 && !d.emulsifierUsed

	options := []uint64{}
	snapshot := vialprofiles.GetSystemVialConfigurationSnapshot()
	for i, p := range snapshot.Profiles {
		if p.VialFluid == machinepb.VialProfile_VIAL_FLUID_DYE_WATER_BASED ||
			(p.VialFluid == machinepb.VialProfile_VIAL_FLUID_EMULSIFIER && emulsifierAllowed) {
			options = append(options, i)
		}
	}

	if len(options) == 0 {
		fmt.Println("ERROR: no system vial profiles")
		return 0
	}
	choiceIndex := d.rand.Intn(len(options))

	choice := options[choiceIndex]
	if snapshot.Profiles[choice].VialFluid == machinepb.VialProfile_VIAL_FLUID_EMULSIFIER {
		d.emulsifierUsed = true
	}
	return choice
}

func (d *autoDecider) decideCollection(predictedState *machinepb.StateReport) *types.CollectionDecision {
	return &types.CollectionDecision{
		VialNo:  int(d.GetRandomVialPos(predictedState)),
		DropsNo: 3,
	}
}

// decideDispense decides a random location from the unit circle
func (d *autoDecider) decideDispense(predictedState *machinepb.StateReport) *types.DispenseDecision {
	x, y := d.sampleRandomUnitCircleCoordinate()
	return &types.DispenseDecision{
		X: float32(x),
		Y: float32(y),
	}
}

func (d *autoDecider) DecideNextAction(predictedState *machinepb.StateReport) (executor.Executor, error) {
	if predictedState.Status <= machinepb.Status_SHUTTING_DOWN {
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

// decideLocationWithinCircle uses rejection sampling to generate coordinates
// within a 1 unit radius circle
func (d *autoDecider) sampleRandomUnitCircleCoordinate() (x, y float64) {
	x = d.rand.Float64()*2 - 1
	y = d.rand.Float64()*2 - 1
	if math.Hypot(x, y) > 1 {
		return d.sampleRandomUnitCircleCoordinate()
	}
	return x, y
}
