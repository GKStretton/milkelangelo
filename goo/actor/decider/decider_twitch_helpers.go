package decider

import (
	"fmt"
	"strings"
	"time"

	"github.com/gkstretton/dark/services/goo/types"
	"github.com/gkstretton/dark/services/goo/vialprofiles"
)

func (d *twitchDecider) subscribeChatVotes(voteType types.VoteType, vialPosToName map[uint64]string) (voteCh chan *types.Vote, unsubscribe func()) {

	msgCh := d.api.SubscribeChat()

	voteCh = make(chan *types.Vote)
	go func() {
		for {
			msg, ok := <-msgCh
			if !ok {
				close(voteCh)
				return
			}
			vote, err := types.TwitchMessageToVote(voteType, msg, vialPosToName)
			if err != nil {
				fmt.Printf("failed to parse vote from %s: %s\n", msg.Message, err)
				d.api.Reply(msg.ID, err.Error())
				continue
			}
			if vote == nil {
				continue
			}
			if vote.Data.LocationVote != nil && vote.Data.LocationVote.N > 2 {
				n := vote.Data.LocationVote.N
				d.api.Reply(msg.ID, fmt.Sprintf("%dD%s", n, strings.Repeat("!?", int(n))))
			}
			voteCh <- vote
		}
	}()

	return voteCh, func() {
		d.api.UnsubscribeChat(msgCh)
	}
}

// if early exit is returned true after a vote, the vote will finish before the timeout
func conductVotingRound(ebsCh, chatCh <-chan *types.Vote, timeout time.Duration, handler func(*types.Vote) (earlyExit bool)) {
	timeoutCh := time.After(timeout)
	for {
		select {
		case <-timeoutCh:
			return
		case vote := <-ebsCh:
			earlyExit := handler(vote)
			if earlyExit {
				return
			}
		case vote := <-chatCh:
			earlyExit := handler(vote)
			if earlyExit {
				return
			}
		}
	}
}

type RunningAverage struct {
	Count   int
	Average float32
}

func (r *RunningAverage) AddNumber(number float32) {
	r.Count++
	r.Average += (number - r.Average) / float32(r.Count)
}

func getVialOptionsAndMap() ([]string, map[uint64]string) {
	vialConfig := vialprofiles.GetSystemVialConfigurationSnapshot()

	options := []string{}
	vialPosToName := map[uint64]string{}
	for posNo, profile := range vialConfig.GetProfiles() {
		vialPosToName[posNo] = profile.Name
		options = append(options, strings.ToLower(profile.Name))
	}
	return options, vialPosToName
}
