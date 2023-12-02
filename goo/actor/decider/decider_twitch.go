package decider

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/dark/services/goo/actor/executor"
	"github.com/gkstretton/dark/services/goo/ebsinterface"
	"github.com/gkstretton/dark/services/goo/twitchapi"
	"github.com/gkstretton/dark/services/goo/types"
)

type twitchDecider struct {
	ebs *ebsinterface.ExtensionSession
	api *twitchapi.TwitchApi
	// if there are no votes, fallback to this one
	fallback Decider
}

func NewTwitchDecider(ebs *ebsinterface.ExtensionSession, twitchApi *twitchapi.TwitchApi) Decider {
	return &twitchDecider{
		ebs: ebs,
		api: twitchApi,
		// todo: change to a more comprehensive auto decider
		fallback: NewMockDecider(),
	}
}

func (d *twitchDecider) DecideCollection(predictedState *machinepb.StateReport) executor.Executor {
	options, vialPosToName := getVialOptionsAndMap()

	// votes from twitch ebs
	ebsCh := d.ebs.SubscribeVotes()
	defer d.ebs.UnsubscribeVotes(ebsCh)

	// votes from twitch chat
	chatVoteCh, unSub := d.subscribeChatVotes(types.VoteTypeCollection, vialPosToName)
	defer unSub()

	d.api.Announce("Taking votes on next collection. Options: "+strings.Join(options, ", "), twitchapi.COLOUR_GREEN)

	// vialPos -> number of votes
	votes := map[uint64]uint64{}

	conductVotingRound(
		ebsCh,
		chatVoteCh,
		time.Duration(30)*time.Second,
		func(vote *types.Vote) bool {
			data := vote.Data.CollectionVote
			votes[data.VialNo]++
			return false
		},
	)

	// build sorted results
	type voteResult struct {
		pos   uint64
		name  string
		count uint64
	}
	sortedResults := []voteResult{}
	for pos, count := range votes {
		sortedResults = append(sortedResults, voteResult{
			pos:   pos,
			name:  vialPosToName[pos],
			count: count,
		})
	}
	slices.SortFunc(sortedResults, func(a, b voteResult) int { return int(a.count) - int(b.count) })

	if len(sortedResults) == 0 {
		d.api.Say("No votes! Choosing at random...")
		return d.fallback.DecideCollection(predictedState)
	}

	d.api.Say(fmt.Sprintf("Vote settled on %s! Results:\n", sortedResults[0].name))
	for _, res := range sortedResults {
		d.api.Say(fmt.Sprintf("    %s: %d", res.name, res.count))
	}

	winnerId := sortedResults[0].pos

	return executor.NewCollectionExecutor(int(winnerId), int(getVialVolume(int(winnerId))))
}

func (d *twitchDecider) DecideDispense(predictedState *machinepb.StateReport) executor.Executor {
	// votes from twitch ebs
	ebsCh := d.ebs.SubscribeVotes()
	defer d.ebs.UnsubscribeVotes(ebsCh)

	// votes from twitch chat
	chatVoteCh, unSub := d.subscribeChatVotes(types.VoteTypeCollection, nil)
	defer unSub()

	e := executor.NewDispenseExecutor(0, 0)
	x := RunningAverage{}
	y := RunningAverage{}

	conductVotingRound(
		ebsCh,
		chatVoteCh,
		time.Duration(30)*time.Second,
		func(vote *types.Vote) bool {
			data := vote.Data.LocationVote
			x.AddNumber(data.X)
			y.AddNumber(data.Y)
			e.X = x.Average
			e.Y = y.Average
			e.Preempt()
			return false
		},
	)

	if x.Count == 0 {
		d.api.Say("No votes! Choosing at random...")
		return d.fallback.DecideDispense(predictedState)
	}

	d.api.Say("Vote settled on average: %.2f, %.2f!")

	return e
}
