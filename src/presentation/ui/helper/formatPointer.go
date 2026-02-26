package uiHelper

import (
	"fmt"

	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

func FormatPointer[ParamType interface{}](pointer *ParamType) string {
	if pointer == nil {
		return "-"
	}

	switch pointerType := any(*pointer).(type) {
	case tkValueObject.UnixTime:
		return pointerType.ReadRfcDate()
	case tkValueObject.Byte:
		return pointerType.StringWithSuffix()
	}

	return fmt.Sprintf("%v", *pointer)
}
