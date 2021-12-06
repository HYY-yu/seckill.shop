package util

import (
	"reflect"
)

func IsNotZero(data interface{}) bool {
	return !IsZero(data)
}

func IsZero(data interface{}) bool {
	if data == nil {
		return true
	}
	if val, ok := data.(reflect.Value); ok {
		return val.IsZero()
	}
	return reflect.ValueOf(data).IsZero()
}
