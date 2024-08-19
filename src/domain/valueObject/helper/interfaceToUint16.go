package voHelper

import (
	"errors"
	"reflect"
	"strconv"
)

func InterfaceToUint16(input interface{}) (output uint16, err error) {
	var defaultErr error = errors.New("CannotConvertToUint16")

	switch v := input.(type) {
	case string:
		uint64Value, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return 0, defaultErr
		}
		if uint64Value > 65535 {
			return 0, defaultErr
		}
		output = uint16(uint64Value)
	case int, int8, int16, int32, int64:
		intValue := reflect.ValueOf(v).Int()
		if intValue < 0 || intValue > 65535 {
			err = defaultErr
		}
		output = uint16(intValue)
	case uint, uint8, uint16, uint32, uint64:
		uintValue := reflect.ValueOf(v).Uint()
		if uintValue > 65535 {
			err = defaultErr
		}
		output = uint16(uintValue)
	case float32, float64:
		floatValue := reflect.ValueOf(v).Float()
		if floatValue < 0 || floatValue > 65535 {
			err = defaultErr
		}
		output = uint16(floatValue)
	default:
		err = defaultErr
	}

	return output, err
}
