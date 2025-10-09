//nolint:revive
package utils

func Ptr[T any](v T) *T {
	return &v
}

func PtrToVal[T any](ptr *T) T {
	var def T

	if ptr != nil {
		return *ptr
	}

	return def
}

// Opt optional chaining
func Opt[T any](object *T) T {
	if object != nil {
		return *object
	}

	var tmp T
	return tmp
}
