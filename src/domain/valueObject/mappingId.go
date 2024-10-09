package valueObject

import (
	"errors"
	"strconv"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

type MappingId uint64

func NewMappingId(value interface{}) (mappingId MappingId, err error) {
	uintValue, err := voHelper.InterfaceToUint64(value)
	if err != nil {
		return mappingId, errors.New("MappingIdMustBeUint64")
	}

	return MappingId(uintValue), nil
}

func (vo MappingId) Uint64() uint64 {
	return uint64(vo)
}

func (vo MappingId) String() string {
	return strconv.FormatUint(uint64(vo), 10)
}
