package valueObject

import (
	"errors"
	"strings"

	"golang.org/x/exp/slices"
)

type MarketplaceItemType string

var ValidMarketplaceItemTypes = []string{
	"app",
	"framework",
	"stack",
}

func NewMarketplaceItemType(value string) (MarketplaceItemType, error) {
	value = strings.ToLower(value)

	mit := MarketplaceItemType(value)
	if !mit.isValid() {
		return "", errors.New("InvalidMarketplaceItemType")
	}

	return MarketplaceItemType(value), nil
}

func NewMarketplaceItemTypePanic(value string) MarketplaceItemType {
	mit, err := NewMarketplaceItemType(value)
	if err != nil {
		panic(err)
	}

	return mit
}

func (mit MarketplaceItemType) isValid() bool {
	return slices.Contains(ValidMarketplaceItemTypes, string(mit))
}

func (mit MarketplaceItemType) String() string {
	return string(mit)
}
