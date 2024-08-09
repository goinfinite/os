package voHelper

import (
	"errors"
	"reflect"
	"strconv"
)

func InterfaceToInt64(input interface{}) (output int64, err error) {
	defaultErr := errors.New("CannotConvertToInt64")

	switch v := input.(type) {
	case string:
		output, err = strconv.ParseInt(v, 10, 64)
	case int, int8, int16, int32, int64:
		intValue := reflect.ValueOf(v).Int()
		output = int64(intValue)
	case uint, uint8, uint16, uint32, uint64:
		uintValue := reflect.ValueOf(v).Uint()
		output = int64(uintValue)
	case float32, float64:
		floatValue := reflect.ValueOf(v).Float()
		output = int64(floatValue)
	default:
		err = defaultErr
	}

	if err != nil {
		return 0, defaultErr
	}

	return output, nil
}
