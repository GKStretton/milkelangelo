package gooapi

import "github.com/op/go-logging"

var l = logging.MustGetLogger("gooapi")

type GooApi interface {
	CollectFromVial(vial int) error
	Dispense(x, y float32) error
	GoToPosition(x, y float32) error
	SetStateUpdateCallback(callback func(state GooStateUpdate))
}
