package valueObject

import (
	"errors"
	"regexp"
	"strings"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

type MarketplaceItemSlug string

const marketplaceItemSlugRegexExpression = `^[a-z0-9\_\-]{2,64}$`

func NewMarketplaceItemSlug(value interface{}) (
	marketplaceItemSlug MarketplaceItemSlug, err error,
) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return marketplaceItemSlug, errors.New("MarketplaceItemSlugValueMustBeString")
	}
	stringValue = strings.TrimSpace(stringValue)
	stringValue = strings.ToLower(stringValue)

	re := regexp.MustCompile(marketplaceItemSlugRegexExpression)
	if !re.MatchString(stringValue) {
		return marketplaceItemSlug, errors.New("InvalidMarketplaceItemSlug")
	}

	return MarketplaceItemSlug(stringValue), nil
}

func (vo MarketplaceItemSlug) String() string {
	return string(vo)
}
