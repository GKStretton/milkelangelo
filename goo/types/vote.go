package types

import "fmt"

type VoteType string

const (
	VoteTypeLocation   VoteType = "LOCATION"
	VoteTypeCollection VoteType = "COLLECTION"
)

type CollectionVote struct {
	VialNo int
}

type LocationVote struct {
	X float32
	Y float32
}

type VoteDetails struct {
	VoteType       VoteType
	CollectionVote *CollectionVote
	LocationVote   *LocationVote
}

type Vote struct {
	Data          VoteDetails
	OpaqueUserID  string
	IsBroadcaster bool
}

func (v *Vote) String() string {
	d := v.Data
	switch d.VoteType {
	case VoteTypeCollection:
		return fmt.Sprintf("%s (%+v)", d.VoteType, d.CollectionVote)
	case VoteTypeLocation:
		return fmt.Sprintf("%s (%+v)", d.VoteType, d.LocationVote)
	default:
		return fmt.Sprintf("unrecognised type %s", d.VoteType)
	}
}
