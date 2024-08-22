package valueObject

import (
	"errors"
	"slices"
	"strings"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

type MappingMatchPattern string

var validMappingMatchPatterns = []string{
	"begins-with", "contains", "equals", "ends-with",
}

func NewMappingMatchPattern(value interface{}) (
	mappingMatchPattern MappingMatchPattern, err error,
) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return mappingMatchPattern, errors.New("MappingMatchPatternMustBeString")
	}
	stringValue = strings.ToLower(stringValue)
	stringValue = strings.ReplaceAll(stringValue, " ", "-")

	if !slices.Contains(validMappingMatchPatterns, stringValue) {
		return mappingMatchPattern, errors.New("InvalidMappingMatchPattern")
	}

	return MappingMatchPattern(stringValue), nil
}

func (vo MappingMatchPattern) String() string {
	return string(vo)
}
