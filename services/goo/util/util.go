package util

import "os"

// Ptr returns a pointer of any value, useful for getting pointers of primitives
func Ptr[T any](v T) *T {
	return &v
}

// EnvBool returns true if the env var is set to true, otherwise false
func EnvBool(name string) bool {
	return os.Getenv(name) == "true"
}
