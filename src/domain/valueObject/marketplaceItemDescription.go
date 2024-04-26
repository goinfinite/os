package valueObject

import (
	"errors"
)

type MarketplaceItemDescription string

func NewMarketplaceItemDescription(value string) (MarketplaceItemDescription, error) {
	if len(value) < 2 {
		return "", errors.New("MarketplaceItemDescriptionTooSmall")
	}

	if len(value) > 2048 {
		return "", errors.New("MarketplaceItemDescriptionTooBig")
	}

	return MarketplaceItemDescription(value), nil
}

func NewMarketplaceItemDescriptionPanic(value string) MarketplaceItemDescription {
	mid, err := NewMarketplaceItemDescription(value)
	if err != nil {
		panic(err)
	}

	return mid
}

func (mid MarketplaceItemDescription) String() string {
	return string(mid)
}
