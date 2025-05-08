package uiHelper

import (
	"fmt"

	"github.com/goinfinite/os/src/domain/valueObject"
)

func FormatPointer[ParamType interface{}](pointer *ParamType) string {
	if pointer == nil {
		return "-"
	}

	switch pointerType := any(*pointer).(type) {
	case valueObject.UnixTime:
		return pointerType.ReadRfcDate()
	case valueObject.Byte:
		return pointerType.StringWithSuffix()
	}

	return fmt.Sprintf("%v", *pointer)
}
