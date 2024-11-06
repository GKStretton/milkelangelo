package types

type EbsMessageType string

const (
	EbsDispenseRequest   EbsMessageType = "dispense"
	EbsCollectionRequest EbsMessageType = "collection"
	EbsGoToRequest       EbsMessageType = "goto"
)

type EbsMessage struct {
	Type              EbsMessageType
	DispenseRequest   *dispenseRequest
	CollectionRequest *collectionRequest
	GoToRequest       *goToRequest
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
