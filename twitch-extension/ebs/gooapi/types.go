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
