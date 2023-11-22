package decider

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/dark/services/goo/actor/executor"
	"github.com/gkstretton/dark/services/goo/livechat"
	"github.com/gkstretton/dark/services/goo/vialprofiles"
)

type twitchDecider struct {
	api *livechat.TwitchApi
	// if there are no votes, fallback to this one
	fallback Decider
}

func NewTwitchDecider(twitchApi *livechat.TwitchApi) Decider {
	return &twitchDecider{
		api: twitchApi,
		// todo: change to a more comprehensive auto decider
		fallback: NewMockDecider(),
	}
}

// if early exit is returned true after a vote, the vote will finish before the timeout
func conductVote(msgCh chan *livechat.Message, timeout time.Duration, handler func(*livechat.Message) (earlyExit bool)) {
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

	vialConfig := vialprofiles.GetSystemVialConfigurationSnapshot()

	options := []string{}
	idToName := map[uint64]string{}
	for _, profile := range vialConfig.GetProfiles() {
		idToName[profile.Id] = profile.Name
		options = append(options, strings.ToLower(profile.Name))
	}
	d.api.Announce("Taking votes on next collection. Options: "+strings.Join(options, ", "), livechat.COLOUR_GREEN)

	// vialNo -> number of votes
	votes := map[uint64]uint64{}

	conductVote(
		msgCh,
		time.Duration(30)*time.Second,
		func(msg *livechat.Message) bool {
			lowerCase := strings.ToLower(msg.Message)
			for _, profile := range vialConfig.GetProfiles() {

				if strings.Contains(lowerCase, profile.Name) {
					fmt.Printf("'%s' parsed as a vote for '%s' (%d)\n", msg.Message, profile.Name, profile.Id)
					votes[profile.Id]++
					return false
				}
			}

			fmt.Printf("'%s' could not be parsed as a vote.\n", msg.Message)
			return false
		},
	)

	// build sorted results
	type voteResult struct {
		id    uint64
		name  string
		count uint64
	}
	sortedResults := []voteResult{}
	for id, count := range votes {
		sortedResults = append(sortedResults, voteResult{
			id:    id,
			name:  idToName[id],
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

	winnerId := sortedResults[0].id

	return executor.NewCollectionExecutor(int(winnerId), int(getVialVolume(int(winnerId))))
}

func (d *twitchDecider) DecideDispense(predictedState *machinepb.StateReport) executor.Executor {
	msgCh := d.api.SubscribeChat()
	defer d.api.UnsubscribeChat(msgCh)

	e := executor.NewDispenseExecutor(0, 0)
	x := RunningAverage{}
	y := RunningAverage{}

	conductVote(
		msgCh,
		time.Duration(30)*time.Second,
		func(msg *livechat.Message) bool {
			vote, err := parseCoordinates(msg.Message)
			if err != nil {
				d.api.Reply(msg.ID, err.Error())
				return false
			}
			if vote == nil {
				return false
			}
			if vote.n > 2 {
				d.api.Reply(msg.ID, fmt.Sprintf("%dD%s", vote.n, strings.Repeat("!?", vote.n)))
			}
			x.AddNumber(vote.x)
			y.AddNumber(vote.y)
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

type coordinateVote struct {
	// how many numbers they provided
	n int
	x float32
	y float32
}

// parseCoordinates takes a string and attempts to extract two decimal numbers as coordinates (x and y).
// It returns a pointer to a coordinateVote struct and an error.
// The function extracts all decimal numbers from the input and checks if at least two are present.
// It parses the first two numbers as x and y coordinates, ensuring they are within the range [-1, 1].
// If the numbers are out of this range, or if there's a parsing error, an error is returned.
func parseCoordinates(input string) (*coordinateVote, error) {
	re := regexp.MustCompile(`-?\d+\.\d+`)
	matches := re.FindAllString(input, -1)
	n := len(matches)

	if n < 2 {
		return nil, nil
	}

	x, err := strconv.ParseFloat(matches[0], 64)
	if err != nil {
		fmt.Println("error parsing x coordinate", err)
		return nil, fmt.Errorf("error parsing x coordinate '%s'", matches[0])
	}

	if x < -1 || x > 1 {
		return nil, fmt.Errorf("x should be between -1 and 1, %f is not", x)
	}

	y, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		fmt.Println("error parsing y coordinate", err)
		return nil, fmt.Errorf("error parsing y coordinate '%s'", matches[1])
	}

	if y < -1 || y > 1 {
		return nil, fmt.Errorf("y should be between -1 and 1, %f is not", y)
	}

	return &coordinateVote{
		n: n,
		x: float32(x),
		y: float32(y),
	}, nil
}

type RunningAverage struct {
	Count   int
	Average float32
}

func (r *RunningAverage) AddNumber(number float32) {
	r.Count++
	r.Average += (number - r.Average) / float32(r.Count)
}
