package util

// Ptr returns a pointer of any value, useful for getting pointers of primitives
func Ptr[T any](v T) *T {
	return &v
}
