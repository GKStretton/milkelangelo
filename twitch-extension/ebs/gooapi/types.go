package gooapi

import (
	"github.com/gkstretton/study-of-light/twitch-ebs/entities"
)

const (
	dispenseRequestType   messageType = "dispense"
	collectionRequestType messageType = "collection"
	goToRequestType       messageType = "goto"
	stateReportType       messageType = "state"
)

type dispenseRequest struct {
	X float32
	Y float32
}

type collectionRequest struct {
	Id int
}

type goToRequest struct {
	X float32
	Y float32
}

type EbsStateReport struct {
	ConnectedUser *entities.User
}

type Status = string

const (
	StatusUnknown               Status = "unknown"
	GooStatusDecidingCollection Status = "deciding-collection"
	GooStatusDecidingDispense   Status = "deciding-dispense"
)

type GooStateUpdate struct {
	Status Status
	X      float32
	Y      float32

	VialProfiles map[int]*VialProfile

	CollectionState *CollectionState
	DispenseState   *DispenseState
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
