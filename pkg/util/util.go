package util

import (
	"reflect"
)

func IsZero(data interface{}) bool {
	if data == nil {
		return true
	}
	if val, ok := data.(reflect.Value); ok {
		return val.IsZero()
	}
	return reflect.ValueOf(data).IsZero()
}
