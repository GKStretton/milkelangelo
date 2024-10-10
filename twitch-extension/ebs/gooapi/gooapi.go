package gooapi

type GooApi interface {
	CollectFromVial(vial int) error
	Dispense() error
	GoToPosition(x, y float32) error
}
