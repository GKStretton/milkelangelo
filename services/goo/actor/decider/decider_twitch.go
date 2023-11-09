package decider

import (
	"fmt"
	"time"

	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/dark/services/goo/actor/executor"
	"github.com/gkstretton/dark/services/goo/livechat"
)

type twitchDecider struct {
	api *livechat.TwitchApi
}

func NewTwitchDecider(twitchApi *livechat.TwitchApi) Decider {
	return &twitchDecider{api: twitchApi}
}

func conductVote(msgCh chan *livechat.Message, timeout time.Duration, handler func(*livechat.Message) bool) {
	timeoutCh := time.After(timeout)
	for {
		select {
		case <-timeoutCh:
			return
		case msg := <-msgCh:
			if handler(msg) {
				return
			}
		}
	}
}

func (d *twitchDecider) DecideCollection(predictedState *machinepb.StateReport) executor.Executor {
	msgCh := d.api.SubscribeChat()
	defer d.api.UnsubscribeChat(msgCh)

	d.api.Announce("Taking votes on next collection (2-6)", livechat.COLOUR_GREEN)

	vialNo := 2

	conductVote(
		msgCh,
		time.Duration(10)*time.Second,
		func(msg *livechat.Message) bool {
			fmt.Println("vote? ", msg)
			switch msg.Message {
			case "2":
				vialNo = 2
			case "3":
				vialNo = 3
			case "4":
				vialNo = 4
			case "5":
				vialNo = 5
				return true
			case "6":
				vialNo = 6
			}
			return false
		},
	)

	d.api.Say(fmt.Sprintf("settled collection vote on %d", vialNo))

	return executor.NewCollectionExecutor(vialNo, int(getVialVolume(vialNo)))
}

func (d *twitchDecider) DecideDispense(predictedState *machinepb.StateReport) executor.Executor {
	//todo: implement this one, test above (special 5)
	e := executor.NewDispenseExecutor(0, 0)

	type voteStruct struct {
		x float32
		y float32
	}
	voteChan := make(chan voteStruct)
	for vote := range voteChan {
		e.X = vote.x
		e.Y = vote.y
		e.Preempt()
	}

	return e
}
