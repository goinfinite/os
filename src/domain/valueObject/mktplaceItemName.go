package valueObject

import (
	"errors"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type MktplaceItemName string

const mktplaceItemNameRegexExpression = `^[a-z0-9\-]{5,30}$`

func NewMktplaceItemName(value string) (MktplaceItemName, error) {
	value = strings.ToLower(value)

	min := MktplaceItemName(value)
	if !min.isValid() {
		return "", errors.New("InvalidMarketplaceItemName")
	}

	return min, nil
}

func NewMktplaceItemNamePanic(value string) MktplaceItemName {
	min, err := NewMktplaceItemName(value)
	if err != nil {
		panic(err)
	}

	return min
}

func (min MktplaceItemName) isValid() bool {
	mktplaceItemNameCompiledRegex := regexp.MustCompile(
		mktplaceItemNameRegexExpression,
	)
	return mktplaceItemNameCompiledRegex.MatchString(string(min))
}

func (min MktplaceItemName) String() string {
	return string(min)
}

func (minPtr *MktplaceItemName) UnmarshalJSON(value []byte) error {
	valueStr := string(value)
	unquotedValue := strings.Trim(valueStr, "\"")

	min, err := NewMktplaceItemName(unquotedValue)
	if err != nil {
		return err
	}

	*minPtr = min
	return nil
}

func (minPtr *MktplaceItemName) UnmarshalYAML(value *yaml.Node) error {
	var valueStr string
	err := value.Decode(&valueStr)
	if err != nil {
		return err
	}

	min, err := NewMktplaceItemName(valueStr)
	if err != nil {
		return err
	}

	*minPtr = min
	return nil
}
