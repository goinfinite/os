package valueObject

import (
	"errors"
	"regexp"
	"strings"
)

type MktplaceItemName string

const mktplaceItemNameRegexExpression = `^[a-z0-9\-]{5,30}$`

func NewMktplaceItemName(value string) (MktplaceItemName, error) {
	value = strings.ToLower(value)

	mkplaceItemInstallStep := MktplaceItemName(value)
	if !mkplaceItemInstallStep.isValid() {
		return "", errors.New("InvalidMarketplaceItemName")
	}

	return mkplaceItemInstallStep, nil
}

func NewMktplaceItemNamePanic(value string) MktplaceItemName {
	mkplaceItemInstallStep, err := NewMktplaceItemName(value)
	if err != nil {
		panic(err)
	}

	return mkplaceItemInstallStep
}

func (mktplaceItemName MktplaceItemName) isValid() bool {
	mktplaceItemNameCompiledRegex := regexp.MustCompile(
		mktplaceItemNameRegexExpression,
	)
	return mktplaceItemNameCompiledRegex.MatchString(string(mktplaceItemName))
}

func (mktplaceItemName MktplaceItemName) String() string {
	return string(mktplaceItemName)
}

func (mktplaceItemNamePtr *MktplaceItemName) UnmarshalJSON(value []byte) error {
	valueStr := string(value)
	unquotedValue := strings.Trim(valueStr, "\"")

	mktplaceItemName, err := NewMktplaceItemName(unquotedValue)
	if err != nil {
		return err
	}

	*mktplaceItemNamePtr = mktplaceItemName
	return nil
}
