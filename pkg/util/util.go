package util

import (
	"reflect"
)

// IsNotZero just IsZero invert
func IsNotZero(data interface{}) bool {
	return !IsZero(data)
}

// IsZero check data is zero value
// be care of the slice, a slice like []int{} is not zero
// because in Go, slice pointed to the underlying array,
// so slice is not zero.
func IsZero(data interface{}) bool {
	if data == nil {
		return true
	}
	if val, ok := data.(reflect.Value); ok {
		return val.IsZero()
	}
	return reflect.ValueOf(data).IsZero()
}
