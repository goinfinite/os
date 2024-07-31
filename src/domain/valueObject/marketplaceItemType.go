package valueObject

import (
	"errors"
	"slices"
	"strings"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

type MarketplaceItemType string

var validMarketplaceItemTypes = []string{
	"app",
	"framework",
	"stack",
}

func NewMarketplaceItemType(value interface{}) (
	marketplaceItemType MarketplaceItemType, err error,
) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return marketplaceItemType, errors.New("MarketplaceItemTypeMustBeString")
	}
	stringValue = strings.ToLower(stringValue)

	if !slices.Contains(validMarketplaceItemTypes, stringValue) {
		return marketplaceItemType, errors.New("InvalidMarketplaceItemType")
	}

	return MarketplaceItemType(stringValue), nil
}

func (vo MarketplaceItemType) String() string {
	return string(vo)
}
