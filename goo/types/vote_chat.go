package types

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gkstretton/dark/services/goo/twitchapi"
)

func TwitchMessageToVote(voteType VoteType, msg *twitchapi.Message, vialPosToName map[uint64]string) (*Vote, error) {
	var voteDetails VoteDetails
	switch voteType {
	case VoteTypeLocation:
		data, err := parseCoordinates(msg.Message)
		if err != nil {
			return nil, err
		}
		if data == nil {
			return nil, nil
		}
		voteDetails = VoteDetails{
			VoteType:     voteType,
			LocationVote: data,
		}
	case VoteTypeCollection:
		data, err := parseCollection(msg.Message, vialPosToName)
		if err != nil {
			return nil, err
		}
		if data == nil {
			return nil, nil
		}
		voteDetails = VoteDetails{
			VoteType:       voteType,
			CollectionVote: data,
		}
	}
	return &Vote{
		Data:          voteDetails,
		OpaqueUserID:  msg.User.ID,
		IsBroadcaster: msg.IsSelf(),
	}, nil
}

func parseCollection(input string, vialPosToName map[uint64]string) (*CollectionVote, error) {
	lowerCase := strings.ToLower(input)
	for pos, name := range vialPosToName {

		if strings.Contains(lowerCase, name) {
			fmt.Printf("'%s' parsed as a vote for '%s' (pos %d)\n", input, name, pos)
			return &CollectionVote{
				VialNo: pos,
			}, nil
		}
	}

	return nil, nil
}

// parseCoordinates takes a string and attempts to extract two decimal numbers as coordinates (x and y).
// It returns a pointer to a coordinateVote struct and an error.
// The function extracts all decimal numbers from the input and checks if at least two are present.
// It parses the first two numbers as x and y coordinates, ensuring they are within the range [-1, 1].
// If the numbers are out of this range, or if there's a parsing error, an error is returned.
func parseCoordinates(input string) (*LocationVote, error) {
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

	return &LocationVote{
		N: uint64(n),
		X: float32(x),
		Y: float32(y),
	}, nil
}
