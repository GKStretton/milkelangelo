package decider

import (
	"fmt"
	"strings"
	"time"

	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/dark/services/goo/actor/executor"
	"github.com/gkstretton/dark/services/goo/ebsinterface"
	"github.com/gkstretton/dark/services/goo/twitchapi"
	"github.com/gkstretton/dark/services/goo/types"
	"github.com/gkstretton/dark/services/goo/vialprofiles"
)

type twitchDecider struct {
	ebs *ebsinterface.ExtensionSession
	api *twitchapi.TwitchApi
	// if there are no votes, fallback to this one
	fallback      Decider
	votingTimeout time.Duration
}

func NewTwitchDecider(ebs *ebsinterface.ExtensionSession, twitchApi *twitchapi.TwitchApi, votingTimeout time.Duration, fallback Decider) Decider {
	return &twitchDecider{
		ebs:           ebs,
		api:           twitchApi,
		fallback:      fallback,
		votingTimeout: votingTimeout,
	}
}

func (d *twitchDecider) DecideCollection(predictedState *machinepb.StateReport) *types.CollectionDecision {
	options, vialPosToName := vialprofiles.GetVialOptionsAndMap()

	// votes from twitch ebs
	ebsCh := d.ebs.SubscribeVotes()
	defer d.ebs.UnsubscribeVotes(ebsCh)

	// votes from twitch chat
	chatVoteCh, unSub := d.api.SubscribeChatVotes(types.VoteTypeCollection)
	defer unSub()

	d.api.Announce("Taking votes on next collection. Options: "+strings.Join(options, ", "), twitchapi.COLOUR_GREEN)

	computerVote := d.fallback.DecideCollection(predictedState)

	// vialPos -> number of votes
	votes := map[uint64]uint64{}
	// total number of votes
	var n int

	getVoteStatus := func() *types.VoteStatus {
		return &types.VoteStatus{
			VoteType: types.VoteTypeCollection,
			CollectionVoteStatus: &types.CollectionVoteStatus{
				TotalVotes:    n,
				VoteCounts:    votes,
				VialPosToName: vialPosToName,
				ComputerVote:  computerVote.VialNo,
			},
		}
	}
	d.ebs.UpdateCurrentVoteStatus(getVoteStatus())
	d.ebs.ManualTriggerBroadcast()

	conductVotingRound(
		types.VoteTypeCollection,
		ebsCh,
		chatVoteCh,
		d.votingTimeout,
		func(vote *types.Vote) bool {
			data := vote.Data.CollectionVote
			// ensure vial is valid
			if _, ok := vialPosToName[data.VialNo]; !ok {
				return false
			}

			n++
			votes[data.VialNo]++

			d.ebs.UpdateCurrentVoteStatus(getVoteStatus())

			return false
		},
	)

	d.ebs.UpdateCurrentVoteStatus(nil)
	d.ebs.UpdatePreviousVoteResult(getVoteStatus())
	d.ebs.ManualTriggerBroadcast()

	// build sorted results
	sortedResults := calculateCollectionVoteResults(votes, vialPosToName)

	if len(sortedResults) == 0 {
		msg := fmt.Sprintf("No votes! Choosing %s", vialPosToName[uint64(computerVote.VialNo)])
		d.api.Say(msg)
		l.Println(msg)
		return computerVote
	}

	msg := fmt.Sprintf("Vote settled on %s! (%d votes)", sortedResults[0].name, sortedResults[0].count)
	d.api.Say(msg)
	l.Println(msg)
	for _, res := range sortedResults {
		l.Printf("    %s: %d\n", res.name, res.count)
	}

	winnerId := sortedResults[0].pos

	return &types.CollectionDecision{
		VialNo:  int(winnerId),
		DropsNo: 4, // todo: vote on this?
	}
}

func (d *twitchDecider) DecideDispense(predictedState *machinepb.StateReport) *types.DispenseDecision {
	// votes from twitch ebs
	ebsCh := d.ebs.SubscribeVotes()
	defer d.ebs.UnsubscribeVotes(ebsCh)

	// votes from twitch chat
	chatVoteCh, unSub := d.api.SubscribeChatVotes(types.VoteTypeLocation)
	defer unSub()

	d.api.Announce("Taking votes on next dispense. Chat format 'x, y'", twitchapi.COLOUR_GREEN)

	computerVote := d.fallback.DecideDispense(predictedState)

	e := executor.NewDispenseExecutor(computerVote)
	e.Preempt()

	x := RunningAverage{}
	y := RunningAverage{}

	getVoteStatus := func() *types.VoteStatus {
		return &types.VoteStatus{
			VoteType: types.VoteTypeLocation,
			LocationVoteStatus: &types.LocationVoteStatus{
				TotalVotes: x.Count,
				AverageVote: types.LocationVote{
					X: x.Average,
					Y: y.Average,
				},
				ComputerVote: types.LocationVote{
					X: computerVote.X,
					Y: computerVote.Y,
				},
			},
		}
	}
	d.ebs.UpdateCurrentVoteStatus(getVoteStatus())
	d.ebs.ManualTriggerBroadcast()

	conductVotingRound(
		types.VoteTypeLocation,
		ebsCh,
		chatVoteCh,
		d.votingTimeout,
		func(vote *types.Vote) bool {
			data := vote.Data.LocationVote
			x.AddNumber(data.X)
			y.AddNumber(data.Y)

			d.ebs.UpdateCurrentVoteStatus(getVoteStatus())

			// move to current vote
			e.X = x.Average
			e.Y = y.Average
			e.Preempt()

			return false
		},
	)

	d.ebs.UpdateCurrentVoteStatus(nil)
	d.ebs.UpdatePreviousVoteResult(getVoteStatus())
	d.ebs.ManualTriggerBroadcast()

	if x.Count == 0 {
		msg := fmt.Sprintf("No votes! Choosing (%.2f, %.2f)", computerVote.X, computerVote.Y)
		d.api.Say(msg)
		l.Println(msg)
		return computerVote
	}

	msg := fmt.Sprintf("Choosing vote average (%.2f, %.2f)", x.Average, y.Average)
	d.api.Say(msg)
	l.Println(msg)

	return &types.DispenseDecision{
		X: x.Average,
		Y: y.Average,
	}
}
