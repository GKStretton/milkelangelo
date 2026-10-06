package types

// GooStateUpdater receives changes to the state that is served to public api
// clients.
type GooStateUpdater interface {
	UpdateState(func(state *GooState))
}

type GooStatus = string

const (
	GooStatusUnknown  GooStatus = "unknown"
	GooStatusSleeping GooStatus = "sleeping"
)

type GooState struct {
	Status       GooStatus
	X            float32
	Y            float32
	VialProfiles map[int]*VialProfile

	CollectionState *CollectionState
	DispenseState   *DispenseState

	WaitingForCollection bool
	WaitingForDispense   bool
	ControlEnabled       bool
}

type CollectionState struct {
	VialNumber int
	VolumeUl   float32
	Completed  bool
}

type DispenseState struct {
	VialNumber        int
	VolumeRemainingUl float32
	Completed         bool
}

type VialProfile struct {
	ID           int
	Name         string
	Colour       string
	DropVolumeUl float32
}
