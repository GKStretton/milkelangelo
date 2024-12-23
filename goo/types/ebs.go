package types

type GooStatus = string

const (
	GooStatusUnknown GooStatus = "unknown"
)

type GooState struct {
	Status       GooStatus
	X            float32
	Y            float32
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

type EbsMessageType string

const (
	EbsDispenseRequest   EbsMessageType = "dispense"
	EbsCollectionRequest EbsMessageType = "collection"
	EbsGoToRequest       EbsMessageType = "goto"
	EbsStateReportType   EbsMessageType = "state"
)

type EbsMessage struct {
	Type              EbsMessageType
	DispenseRequest   *dispenseRequest
	CollectionRequest *collectionRequest
	GoToRequest       *goToRequest
	StateReport       *EbsStateReport
}

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
	ConnectedUser *EbsUser
}

type EbsUser struct {
	OUID string
}
