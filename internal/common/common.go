package common

// Make a pointer to any
func MakePtr[T any](v T) *T {
	return &v
}

type RequestID struct{}
type RemoteAddr struct{}
