package voHelper

import (
	"errors"
	"reflect"
	"strconv"
)

func InterfaceToUint(input interface{}) (uint, error) {
	var output uint
	var err error
	var defaultErr error = errors.New("InvalidInput")
	switch v := input.(type) {
	case string:
		uint64Value, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return 0, defaultErr
		}
		if uint64Value > 4294967295 {
			return 0, defaultErr
		}
		output = uint(uint64Value)
	case int, int8, int16, int32, int64:
		intValue := reflect.ValueOf(v).Int()
		if intValue < 0 || intValue > 4294967295 {
			err = defaultErr
		}
		output = uint(intValue)
	case uint, uint8, uint16, uint32, uint64:
		uintValue := uint(reflect.ValueOf(v).Uint())
		if uintValue > 4294967295 {
			err = defaultErr
		}
		output = uintValue
	case float32, float64:
		floatValue := reflect.ValueOf(v).Float()
		if floatValue < 0 || floatValue > 4294967295 {
			err = defaultErr
		}
		output = uint(floatValue)
	default:
		err = defaultErr
	}

	if err != nil {
		return 0, defaultErr
	}

	return output, nil
}
