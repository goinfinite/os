package valueObject

import (
	"errors"
	"regexp"
	"strings"
)

type MarketplaceItemName string

const marketplaceItemNameRegexExpression = `^\p{L}[\p{L}\'\ \-]{3,30}$`

func NewMarketplaceItemName(value string) (MarketplaceItemName, error) {
	value = strings.ToLower(value)

	min := MarketplaceItemName(value)
	if !min.isValid() {
		return "", errors.New("InvalidMarketplaceItemName")
	}

	return min, nil
}

func NewMarketplaceItemNamePanic(value string) MarketplaceItemName {
	min, err := NewMarketplaceItemName(value)
	if err != nil {
		panic(err)
	}

	return min
}

func (min MarketplaceItemName) isValid() bool {
	marketplaceItemNameCompiledRegex := regexp.MustCompile(
		marketplaceItemNameRegexExpression,
	)
	return marketplaceItemNameCompiledRegex.MatchString(string(min))
}

func (min MarketplaceItemName) String() string {
	return string(min)
}
