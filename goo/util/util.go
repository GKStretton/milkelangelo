package util

import (
	"math"
	"math/rand"
	"os"
)

// Ptr returns a pointer of any value, useful for getting pointers of primitives
func Ptr[T any](v T) *T {
	return &v
}

// EnvBool returns true if the env var is set to true, otherwise false
func EnvBool(name string) bool {
	return os.Getenv(name) == "true"
}

// decideLocationWithinCircle uses rejection sampling to generate coordinates
// within a 1 unit radius circle
func SampleRandomUnitCircleCoordinate() (x, y float64) {
	x = rand.Float64()*2 - 1
	y = rand.Float64()*2 - 1
	if math.Hypot(x, y) > 1 {
		return SampleRandomUnitCircleCoordinate()
	}
	return x, y
}
