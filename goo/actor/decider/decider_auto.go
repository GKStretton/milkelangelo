package decider

import (
	"time"

	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/dark/services/goo/actor/executor"
	"github.com/gkstretton/dark/services/goo/types"
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

func (d *autoDecider) decideCollection(predictedState *machinepb.StateReport) *types.CollectionDecision {
	// return &types.CollectionDecision{
	// 	VialNo: 1,
	// 	DropsNo: 1,
	// }
	return nil
}

func (d *autoDecider) decideDispense(predictedState *machinepb.StateReport) *types.DispenseDecision {
	// return &types.DispenseDecision{
	// 	X: 0,
	// 	Y: 0,
	// }
	return nil
}

func (d *autoDecider) DecideNextAction(predictedState *machinepb.StateReport) executor.Executor {
	if time.Now().After(d.endTime) {
		l.Println("endTime reached on auto decider, deciding nil.")
		return nil
	}
	if predictedState.Status == machinepb.Status_SLEEPING {
		l.Println("invalid state for actor, decided nil.")
		return nil
	}
	if predictedState.PipetteState.Spent {
		l.Println("collection is next, launching decider...")
		decision := d.decideCollection(predictedState)
		return executor.NewCollectionExecutor(decision)
	}
	l.Println("dispense is next, launching decider...")
	decision := d.decideDispense(predictedState)
	return executor.NewDispenseExecutor(decision)
}
