package valueObject

import tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"

func NewUnixFileName(
	value interface{},
) (tkValueObject.UnixFileName, error) {
	return tkValueObject.NewUnixFileName(value, false)
}
