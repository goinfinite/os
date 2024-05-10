package valueObject

import (
	"errors"
	"slices"
	"strings"
)

type MappingMatchPattern string

var ValidMappingMatchPatterns = []string{
	"begins-with",
	"contains",
	"equals",
	"ends-with",
}

func NewMappingMatchPattern(value string) (MappingMatchPattern, error) {
	value = strings.ToLower(value)
	if !slices.Contains(ValidMappingMatchPatterns, value) {
		return "", errors.New("InvalidMappingMatchPattern")
	}
	return MappingMatchPattern(value), nil
}

func NewMappingMatchPatternPanic(value string) MappingMatchPattern {
	mmp, err := NewMappingMatchPattern(value)
	if err != nil {
		panic(err)
	}
	return mmp
}

func (mmp MappingMatchPattern) String() string {
	return string(mmp)
}
