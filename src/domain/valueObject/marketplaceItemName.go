package valueObject

import (
	"errors"
	"regexp"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

type MarketplaceItemName string

var marketplaceItemNameRegex = regexp.MustCompile(`^[A-z0-9][\p{L}0-9\'\ \-]{1,30}$`)

func NewMarketplaceItemName(value any) (
	marketplaceItemName MarketplaceItemName, err error,
) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
	if err != nil {
		return marketplaceItemName, errors.New("MarketplaceItemNameMustBeString")
	}

	if !marketplaceItemNameRegex.MatchString(stringValue) {
		return marketplaceItemName, errors.New("InvalidMarketplaceItemName")
	}

	return MarketplaceItemName(stringValue), nil
}

func (vo MarketplaceItemName) String() string {
	return string(vo)
}
