package valueObject

import (
	"errors"
	"regexp"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

var mappingSecurityRuleDescriptionRegex = regexp.MustCompile(`^[^\r\n\t\x00-\x1F\x7F]{0,1000}$`)

type MappingSecurityRuleDescription string

func NewMappingSecurityRuleDescription(value any) (
	mappingSecurityRuleDescription MappingSecurityRuleDescription,
	err error,
) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
	if err != nil {
		return mappingSecurityRuleDescription, errors.New("MappingSecurityRuleDescriptionMustBeString")
	}

	if !mappingSecurityRuleDescriptionRegex.MatchString(stringValue) {
		return mappingSecurityRuleDescription, errors.New("InvalidMappingSecurityRuleDescription")
	}

	return MappingSecurityRuleDescription(stringValue), nil
}

func (vo MappingSecurityRuleDescription) String() string {
	return string(vo)
}
