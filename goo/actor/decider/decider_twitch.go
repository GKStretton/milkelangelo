package decider

import (
	"fmt"
	"strings"
	"time"

	"github.com/gkstretton/asol-protos/go/machinepb"
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

func (d *twitchDecider) DecideCollection(predictedState *machinepb.StateReport) *CollectionDecision {
	options, vialPosToName := vialprofiles.GetVialOptionsAndMap()

	// votes from twitch ebs
	ebsCh := d.ebs.SubscribeVotes()
	defer d.ebs.UnsubscribeVotes(ebsCh)

	// votes from twitch chat
	chatVoteCh, unSub := d.api.SubscribeChatVotes(types.VoteTypeCollection)
	defer unSub()

	d.api.Announce("Taking votes on next collection. Options: "+strings.Join(options, ", "), twitchapi.COLOUR_GREEN)

	// vialPos -> number of votes
	votes := map[uint64]uint64{}
	// total number of votes
	var n int

	d.ebs.UpdateCurrentVoteStatus(&types.VoteStatus{
		VoteType: types.VoteTypeCollection,
		CollectionVoteStatus: &types.CollectionVoteStatus{
			TotalVotes:    n,
			VoteCounts:    votes,
			VialPosToName: vialPosToName,
		},
	})
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

			d.ebs.UpdateCurrentVoteStatus(&types.VoteStatus{
				VoteType: types.VoteTypeCollection,
				CollectionVoteStatus: &types.CollectionVoteStatus{
					TotalVotes:    n,
					VoteCounts:    votes,
					VialPosToName: vialPosToName,
				},
			})

			return false
		},
	)

	d.ebs.UpdateCurrentVoteStatus(nil)
	d.ebs.UpdatePreviousVoteResult(&types.VoteStatus{
		VoteType: types.VoteTypeCollection,
		CollectionVoteStatus: &types.CollectionVoteStatus{
			TotalVotes:    n,
			VoteCounts:    votes,
			VialPosToName: vialPosToName,
		},
	})
	d.ebs.ManualTriggerBroadcast()

	// build sorted results
	sortedResults := calculateCollectionVoteResults(votes, vialPosToName)

	if len(sortedResults) == 0 {
		d.api.Say("No votes! Choosing at random...")
		l.Println("fallback decider for collection...")
		return d.fallback.DecideCollection(predictedState)
	}

	msg := fmt.Sprintf("Vote settled on %s! (%d votes)", sortedResults[0].name, sortedResults[0].count)
	d.api.Say(msg)
	l.Println(msg)
	for _, res := range sortedResults {
		l.Printf("    %s: %d\n", res.name, res.count)
	}

	winnerId := sortedResults[0].pos

	// return executor.NewCollectionExecutor(int(winnerId), int(getVialVolume(int(winnerId))))
}

func (d *twitchDecider) DecideDispense(predictedState *machinepb.StateReport) *DispenseDecision {
	// votes from twitch ebs
	ebsCh := d.ebs.SubscribeVotes()
	defer d.ebs.UnsubscribeVotes(ebsCh)

	// votes from twitch chat
	chatVoteCh, unSub := d.api.SubscribeChatVotes(types.VoteTypeLocation)
	defer unSub()

	d.api.Announce("Taking votes on next dispense. Chat format 'x, y'", twitchapi.COLOUR_GREEN)

	computerVote := d.fallback.DecideDispense(predictedState)

	// e := executor.NewDispenseExecutor(0, 0)
	x := RunningAverage{}
	y := RunningAverage{}

	d.ebs.UpdateCurrentVoteStatus(&types.VoteStatus{
		VoteType: types.VoteTypeLocation,
		LocationVoteStatus: &types.LocationVoteStatus{
			TotalVotes: x.Count,
			XAvg:       x.Average,
			YAvg:       y.Average,
		},
	})
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
			e.X = x.Average
			e.Y = y.Average
			e.Preempt()

			d.ebs.UpdateCurrentVoteStatus(&types.VoteStatus{
				VoteType: types.VoteTypeLocation,
				LocationVoteStatus: &types.LocationVoteStatus{
					TotalVotes: x.Count,
					XAvg:       x.Average,
					YAvg:       y.Average,
				},
			})

			return false
		},
	)

	d.ebs.UpdateCurrentVoteStatus(nil)
	d.ebs.UpdatePreviousVoteResult(&types.VoteStatus{
		VoteType: types.VoteTypeLocation,
		LocationVoteStatus: &types.LocationVoteStatus{
			TotalVotes: x.Count,
			XAvg:       x.Average,
			YAvg:       y.Average,
		},
	})
	d.ebs.ManualTriggerBroadcast()

	if x.Count == 0 {
		d.api.Say("No votes! Choosing at random...")
		l.Println("fallback decider for location...")
		return d.fallback.DecideDispense(predictedState)
	}

	msg := fmt.Sprintf("Vote settled on average: %.2f, %.2f!", x.Average, y.Average)
	d.api.Say(msg)
	l.Println(msg)

	return e
}
