package valueObject

import (
	"errors"
	"strings"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

type MappingMatchPattern string

const (
	MappingMatchPatternBeginsWith MappingMatchPattern = "begins-with"
	MappingMatchPatternContains   MappingMatchPattern = "contains"
	MappingMatchPatternEquals     MappingMatchPattern = "equals"
	MappingMatchPatternEndsWith   MappingMatchPattern = "ends-with"
)

func NewMappingMatchPattern(value interface{}) (
	matchPattern MappingMatchPattern, err error,
) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return matchPattern, errors.New("MappingMatchPatternMustBeString")
	}
	stringValue = strings.ToLower(stringValue)
	stringValue = strings.ReplaceAll(stringValue, " ", "-")

	stringValueVo := MappingMatchPattern(stringValue)
	switch stringValueVo {
	case MappingMatchPatternBeginsWith, MappingMatchPatternContains,
		MappingMatchPatternEquals, MappingMatchPatternEndsWith:
		return stringValueVo, nil
	default:
		return matchPattern, errors.New("InvalidMappingMatchPattern")
	}
}

func (vo MappingMatchPattern) String() string {
	return string(vo)
}
