package valueObject

import (
	"errors"
	"strings"

	"golang.org/x/exp/slices"
)

type MappingMatchPattern string

var ValidMappingMatchPatterns = []string{
	"contains",
	"equals",
	"beginsWith",
	"endsWith",
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
