package util

import "os"

// Ptr returns a pointer of any value, useful for getting pointers of primitives
func Ptr[T any](v T) *T {
	return &v
}

// EnvPresent returns true if the env var is set to something, otherwise false
func EnvPresent(name string) bool {
	return os.Getenv(name) != ""
}
