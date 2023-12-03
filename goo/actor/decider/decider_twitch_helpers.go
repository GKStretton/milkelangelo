package decider

import (
	"slices"
	"time"

	"github.com/gkstretton/dark/services/goo/types"
)

// if early exit is returned true after a vote, the vote will finish before the timeout
func conductVotingRound(t types.VoteType, ebsCh, chatCh <-chan *types.Vote, timeout time.Duration, handler func(*types.Vote) (earlyExit bool)) {
	timeoutCh := time.After(timeout)
	for {
		select {
		case <-timeoutCh:
			l.Println("vote timed out")
			return
		case vote := <-ebsCh:
			if vote.Data.VoteType != t {
				continue
			}
			l.Println("got ebs vote")
			earlyExit := handler(vote)
			if earlyExit {
				return
			}
		case vote := <-chatCh:
			if vote.Data.VoteType != t {
				continue
			}
			l.Println("got chat vote")
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

type collectionVoteResult struct {
	pos   uint64
	name  string
	count uint64
}

func calculateCollectionVoteResults(votes map[uint64]uint64, vialPosToName map[uint64]string) []collectionVoteResult {
	sortedResults := []collectionVoteResult{}
	for pos, count := range votes {
		sortedResults = append(sortedResults, collectionVoteResult{
			pos:   pos,
			name:  vialPosToName[pos],
			count: count,
		})
	}
	slices.SortFunc(sortedResults, func(a, b collectionVoteResult) int { return int(b.count) - int(a.count) })
	return sortedResults
}
