package valueObject

import (
	"errors"
	"regexp"
	"strings"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

type MarketplaceItemSlug string

var marketplaceItemSlugRegex = regexp.MustCompile(`^[a-z0-9\_\-]{2,64}$`)

func NewMarketplaceItemSlug(value any) (
	marketplaceItemSlug MarketplaceItemSlug, err error,
) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
	if err != nil {
		return marketplaceItemSlug, errors.New("MarketplaceItemSlugValueMustBeString")
	}
	stringValue = strings.ToLower(stringValue)

	if !marketplaceItemSlugRegex.MatchString(stringValue) {
		return marketplaceItemSlug, errors.New("InvalidMarketplaceItemSlug")
	}

	return MarketplaceItemSlug(stringValue), nil
}

func (vo MarketplaceItemSlug) String() string {
	return string(vo)
}
