package valueObject

import (
	"errors"
	"regexp"
	"strings"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

var mappingPathRegex = regexp.MustCompile(`^[^\s<>;'":#{}?\[\]]{1,512}$`)

type MappingPath string

func NewMappingPath(value any) (mappingPath MappingPath, err error) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
	if err != nil {
		return mappingPath, errors.New("MappingPathMustBeString")
	}

	hasLeadingSlash := strings.HasPrefix(stringValue, "/")
	if !hasLeadingSlash {
		stringValue = "/" + stringValue
	}

	if !mappingPathRegex.MatchString(stringValue) {
		return mappingPath, errors.New("InvalidMappingPath")
	}

	return MappingPath(stringValue), nil
}

func (vo MappingPath) String() string {
	return string(vo)
}
