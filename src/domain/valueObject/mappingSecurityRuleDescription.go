package valueObject

import (
	"errors"
	"regexp"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

const mappingSecurityRuleDescriptionRegex string = `^[^\r\n\t\x00-\x1F\x7F]{0,1000}$`

type MappingSecurityRuleDescription string

func NewMappingSecurityRuleDescription(value interface{}) (
	mappingSecurityRuleDescription MappingSecurityRuleDescription,
	err error,
) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
	if err != nil {
		return mappingSecurityRuleDescription, errors.New("MappingSecurityRuleDescriptionMustBeString")
	}

	re := regexp.MustCompile(mappingSecurityRuleDescriptionRegex)
	if !re.MatchString(stringValue) {
		return mappingSecurityRuleDescription, errors.New("InvalidMappingSecurityRuleDescription")
	}

	return MappingSecurityRuleDescription(stringValue), nil
}

func (vo MappingSecurityRuleDescription) String() string {
	return string(vo)
}
