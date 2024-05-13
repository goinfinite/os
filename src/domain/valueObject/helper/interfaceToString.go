package voHelper

import (
	"errors"
	"reflect"
	"strconv"
)

func InterfaceToString(input interface{}) (string, error) {
	var output string

	switch v := input.(type) {
	case string:
		output = v
	case int, int8, int16, int32, int64:
		intValue := reflect.ValueOf(v).Int()
		output = strconv.FormatInt(intValue, 10)
	case uint, uint8, uint16, uint32, uint64:
		uintValue := reflect.ValueOf(v).Uint()
		output = strconv.FormatUint(uintValue, 10)
	case float32, float64:
		floatValue := reflect.ValueOf(v).Float()
		output = strconv.FormatFloat(floatValue, 'f', -1, 64)
	case bool:
		boolValue := reflect.ValueOf(v).Bool()
		output = strconv.FormatBool(boolValue)
	default:
		return "", errors.New("CannotConvertToString")
	}

	return output, nil
}
