package presenterValueObject

import (
	"errors"
	"slices"
	"strings"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

type MarketplaceListType string

var validMarketplaceListTypes = []string{"installed", "catalog"}

func NewMarketplaceListType(value interface{}) (listType MarketplaceListType, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return listType, errors.New("MarketplaceListTypeMustBeString")
	}
	stringValue = strings.ToLower(stringValue)

	if !slices.Contains(validMarketplaceListTypes, stringValue) {
		return listType, errors.New("InvalidMarketplaceListType")
	}

	return MarketplaceListType(stringValue), nil
}

func (vo MarketplaceListType) String() string {
	return string(vo)
}
