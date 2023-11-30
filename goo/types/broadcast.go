package types

type VoteStatus struct {
	VoteType             VoteType
	CollectionVoteStatus *CollectionVoteStatus
	LocationVoteStatus   *LocationVoteStatus
}

type CollectionVoteStatus struct {
	TotalVotes int
	VoteCounts map[int]int
}

type LocationVoteStatus struct {
	TotalVotes int
	XAvg       float32
	YAvg       float32
}
