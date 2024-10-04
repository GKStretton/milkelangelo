// Define VoteType as an enum
export enum VoteType {
	Location = "LOCATION",
	Collection = "COLLECTION",
}

// Define CollectionVote as an interface (or class if you need methods)
export interface CollectionVote {
	vialNo: number;
}

// Define LocationVote as an interface
export interface LocationVote {
	x: number; // TypeScript uses 'number' for all numeric types
	y: number;
}

// Define VoteDetails as an interface
export interface VoteDetails {
	voteType: VoteType;
	collectionVote?: CollectionVote; // Optional property
	locationVote?: LocationVote; // Optional property
}

export function createCollectionVote(vialNo: number): VoteDetails {
	return {
		voteType: VoteType.Collection,
		collectionVote: { vialNo },
	};
}

export function createLocationVote(x: number, y: number): VoteDetails {
	return {
		voteType: VoteType.Location,
		locationVote: { x, y },
	};
}
