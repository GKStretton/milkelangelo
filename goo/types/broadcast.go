package types

type VoteStatus struct {
	VoteType             VoteType
	CollectionVoteStatus *CollectionVoteStatus
	LocationVoteStatus   *LocationVoteStatus
}

type CollectionVoteStatus struct {
	TotalVotes    int
	VoteCounts    map[uint64]uint64
	VialPosToName map[uint64]string
	ComputerVote  int
}

type LocationVoteStatus struct {
	TotalVotes   int
	AverageVote  LocationVote
	ComputerVote LocationVote
}
