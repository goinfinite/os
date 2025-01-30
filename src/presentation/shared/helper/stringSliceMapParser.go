package sharedHelper

import (
	"log/slog"
	"reflect"
	"strings"
)

func StringSliceValueObjectParser[TypedObject any](
	rawInputValues any,
	valueObjectConstructor func(any) (TypedObject, error),
) []TypedObject {
	resultObjects := make([]TypedObject, 0)

	if rawInputValues == nil {
		return resultObjects
	}

	rawReflectedSlice := make([]interface{}, 0)

	reflectedRawValues := reflect.ValueOf(rawInputValues)
	rawInputValuesKind := reflectedRawValues.Kind()
	switch rawInputValuesKind {
	case reflect.String:
		rawInputValues = strings.Split(reflectedRawValues.String(), ";")
		for _, rawValue := range rawInputValues.([]string) {
			rawReflectedSlice = append(rawReflectedSlice, rawValue)
		}
	case reflect.Slice:
		for valueIndex := 0; valueIndex < reflectedRawValues.Len(); valueIndex++ {
			rawReflectedSlice = append(
				rawReflectedSlice, reflectedRawValues.Index(valueIndex).Interface(),
			)
		}
	default:
		rawReflectedSlice = append(rawReflectedSlice, rawInputValues)
	}

	for _, rawValue := range rawReflectedSlice {
		valueObject, err := valueObjectConstructor(rawValue)
		if err != nil {
			slog.Debug(err.Error(), slog.Any("rawValue", rawValue))
			continue
		}

		resultObjects = append(resultObjects, valueObject)
	}

	return resultObjects
}
