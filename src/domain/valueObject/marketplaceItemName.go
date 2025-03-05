package valueObject

import (
	"errors"
	"regexp"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

type MarketplaceItemName string

const marketplaceItemNameRegexExpression = `^[A-z0-9][\p{L}0-9\'\ \-]{1,30}$`

func NewMarketplaceItemName(value interface{}) (
	marketplaceItemName MarketplaceItemName, err error,
) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return marketplaceItemName, errors.New("MarketplaceItemNameMustBeString")
	}

	re := regexp.MustCompile(marketplaceItemNameRegexExpression)
	if !re.MatchString(stringValue) {
		return marketplaceItemName, errors.New("InvalidMarketplaceItemName")
	}

	return MarketplaceItemName(stringValue), nil
}

func (vo MarketplaceItemName) String() string {
	return string(vo)
}
