export interface Coords {
	x: number;
	y: number;
}

export interface EbsState {
	GooState: GooState;
	ConnectedUser: User;
}

export interface GooState {
	Status: Status;

	X: number;
	Y: number;

	VialProfiles: { [key: number]: VialProfile };

	CollectionState: CollectionState | null;
	DispenseState: DispenseState | null;
}

export interface User {
	OUID: string;
}

export enum Status {
	StatusUnknown = "unknown",
	StatusDecidingCollection = "deciding-collection",
	StatusDecidingDispense = "deciding-dispense",
}

export interface CollectionState {
	VialNumber: number;
	VolumeUl: number;
	Completed: boolean;
}

export interface DispenseState {
	VialNumber: number;
	VolumeRemainingUl: number;
	Completed: boolean;
}

export interface VialProfile {
	ID: number;
	Name: string;
	Colour: string;
	DropVolumeUl: number;
}
