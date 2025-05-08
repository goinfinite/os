package valueObject

import (
	"errors"
	"strconv"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

type MappingSecurityRuleId uint64

func NewMappingSecurityRuleId(value interface{}) (
	mappingSecurityRuleId MappingSecurityRuleId,
	err error,
) {
	uintValue, err := tkVoUtil.InterfaceToUint64(value)
	if err != nil {
		return mappingSecurityRuleId, errors.New("MappingSecurityRuleIdMustBeUint64")
	}

	if uintValue == 0 {
		return mappingSecurityRuleId, errors.New("MappingSecurityRuleIdCannotBeZero")
	}

	return MappingSecurityRuleId(uintValue), nil
}

func (vo MappingSecurityRuleId) Uint64() uint64 {
	return uint64(vo)
}

func (vo MappingSecurityRuleId) String() string {
	return strconv.FormatUint(uint64(vo), 10)
}
