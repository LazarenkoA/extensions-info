//nolint:revive
package utils

import (
	"crypto/md5"
	"encoding/hex"
)

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

func Cast[T any](v interface{}) T {
	var defaultVal T
	if v, ok := v.(T); ok {
		return v
	}

	return defaultVal
}

func Hash(backed []byte) string {
	hash := md5.Sum(backed)
	return hex.EncodeToString(hash[:])
}
