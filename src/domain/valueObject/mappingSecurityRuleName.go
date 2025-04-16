package valueObject

import (
	"errors"
	"regexp"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

const mappingSecurityRuleNameRegex string = `^[a-zA-Z0-9][a-zA-Z0-9\-_ ]{1,512}$`

type MappingSecurityRuleName string

func NewMappingSecurityRuleName(value interface{}) (
	mappingSecurityRuleName MappingSecurityRuleName,
	err error,
) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
	if err != nil {
		return mappingSecurityRuleName, errors.New("MappingSecurityRuleNameMustBeString")
	}

	re := regexp.MustCompile(mappingSecurityRuleNameRegex)
	if !re.MatchString(stringValue) {
		return mappingSecurityRuleName, errors.New("InvalidMappingSecurityRuleName")
	}

	return MappingSecurityRuleName(stringValue), nil
}

func (vo MappingSecurityRuleName) String() string {
	return string(vo)
}
