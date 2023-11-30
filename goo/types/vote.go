package types

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
