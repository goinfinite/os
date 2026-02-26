package valueObject

import (
	"errors"
	"slices"
	"strings"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

type MarketplaceItemType string

var validMarketplaceItemTypes = []string{
	"app", "framework", "stack",
}

func NewMarketplaceItemType(value interface{}) (
	marketplaceItemType MarketplaceItemType, err error,
) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
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
