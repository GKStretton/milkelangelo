// Mirrors goo/types/public.go and the state sent by goo/publicapi

export interface RemoteState {
	GooState: GooState;
	// true while any client holds control, including this one
	Claimed: boolean;
}

export interface GooState {
	Status: "unknown" | "sleeping";
	// pipette target, within the unit circle
	X: number;
	Y: number;
	// keyed by vial position; null when a position has no valid profile
	VialProfiles: Record<string, VialProfile | null> | null;
	CollectionState: CollectionState | null;
	DispenseState: DispenseState | null;
	WaitingForCollection: boolean;
	WaitingForDispense: boolean;
	ControlEnabled: boolean;
}

export interface CollectionState {
	VialNumber: number;
	VolumeUl: number;
	Completed: boolean;
}

export interface DispenseState {
	VialNumber: number;
	VolumeRemainingUl: number;
	// true when the pipette is empty
	Completed: boolean;
}

export interface VialProfile {
	ID: number;
	Name: string;
	Colour: string;
	DropVolumeUl: number;
}
