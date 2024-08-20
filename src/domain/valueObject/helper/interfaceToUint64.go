package voHelper

import (
	"errors"
	"reflect"
	"strconv"
)

func InterfaceToUint64(input interface{}) (output uint64, err error) {
	defaultErr := errors.New("CannotConvertToUint64")

	switch v := input.(type) {
	case string:
		output, err = strconv.ParseUint(v, 10, 64)
	case int, int8, int16, int32, int64:
		intValue := reflect.ValueOf(v).Int()
		if intValue < 0 {
			err = defaultErr
		}
		output = uint64(intValue)
	case uint, uint8, uint16, uint32, uint64:
		output = uint64(reflect.ValueOf(v).Uint())
	case float32, float64:
		floatValue := reflect.ValueOf(v).Float()
		if floatValue < 0 {
			err = defaultErr
		}
		output = uint64(floatValue)
	default:
		err = defaultErr
	}

	if err != nil {
		return 0, defaultErr
	}

	return output, nil
}
