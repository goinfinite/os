package valueObject

import (
	"errors"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

type MarketplaceItemDescription string

func NewMarketplaceItemDescription(value interface{}) (
	marketplaceItemDescription MarketplaceItemDescription, err error,
) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
	if err != nil {
		return marketplaceItemDescription, errors.New(
			"MarketplaceItemDescriptionMustBeString",
		)
	}

	if len(stringValue) < 2 {
		return marketplaceItemDescription, errors.New(
			"MarketplaceItemDescriptionTooSmall",
		)
	}

	if len(stringValue) > 2048 {
		return marketplaceItemDescription, errors.New(
			"MarketplaceItemDescriptionTooBig",
		)
	}

	return MarketplaceItemDescription(stringValue), nil
}

func (vo MarketplaceItemDescription) String() string {
	return string(vo)
}
