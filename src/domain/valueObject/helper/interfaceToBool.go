package voHelper

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

func InterfaceToBool(input interface{}) (output bool, err error) {
	switch v := input.(type) {
	case bool:
		return v, nil
	case string:
		stringValue := strings.TrimSpace(v)
		stringValue = strings.ToLower(stringValue)
		if stringValue == "on" {
			return true, nil
		}

		if stringValue == "off" {
			return false, nil
		}

		return strconv.ParseBool(stringValue)
	case int, int8, int16, int32, int64:
		intValue := reflect.ValueOf(v).Int()
		return intValue != 0, nil
	case uint, uint8, uint16, uint32, uint64:
		uintValue := reflect.ValueOf(v).Uint()
		return uintValue != 0, nil
	case float32, float64:
		floatValue := reflect.ValueOf(v).Float()
		return floatValue != 0, nil
	default:
		return false, errors.New("CannotConvertToBool")
	}
}
