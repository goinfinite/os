package valueObject

import (
	"errors"
	"regexp"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

var mappingSecurityRuleNameRegex = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9\-_ ]{1,512}$`)

type MappingSecurityRuleName string

func NewMappingSecurityRuleName(value any) (
	mappingSecurityRuleName MappingSecurityRuleName,
	err error,
) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
	if err != nil {
		return mappingSecurityRuleName, errors.New("MappingSecurityRuleNameMustBeString")
	}

	if !mappingSecurityRuleNameRegex.MatchString(stringValue) {
		return mappingSecurityRuleName, errors.New("InvalidMappingSecurityRuleName")
	}

	return MappingSecurityRuleName(stringValue), nil
}

func (vo MappingSecurityRuleName) String() string {
	return string(vo)
}
