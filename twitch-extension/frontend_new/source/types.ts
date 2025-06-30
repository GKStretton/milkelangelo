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

	WaitingForCollection: boolean;
	WaitingForDispense: boolean;
	ActorRunning: boolean;
}

export interface User {
	OUID: string;
}

export enum Status {
	StatusUnknown = "unknown",
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
