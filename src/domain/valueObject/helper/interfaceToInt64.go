package voHelper

import (
	"errors"
	"reflect"
	"strconv"
)

func InterfaceToInt64(input interface{}) (output int64, err error) {
	switch v := input.(type) {
	case string:
		return strconv.ParseInt(v, 10, 64)
	case int, int8, int16, int32, int64:
		intValue := reflect.ValueOf(v).Int()
		return int64(intValue), nil
	case uint, uint8, uint16, uint32, uint64:
		uintValue := reflect.ValueOf(v).Uint()
		return int64(uintValue), nil
	case float32, float64:
		floatValue := reflect.ValueOf(v).Float()
		return int64(floatValue), nil
	default:
		return 0, errors.New("CannotConvertToInt64")
	}
}
