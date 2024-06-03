package valueObject

import (
	"errors"
	"regexp"
	"strings"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

type MarketplaceItemSlug string

const marketplaceItemSlugRegexExpression = `^[a-z0-9\_\-]{2,64}$`

func NewMarketplaceItemSlug(value interface{}) (MarketplaceItemSlug, error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return "", errors.New("MarketplaceItemSlugValueMustBeString")
	}

	stringValue = strings.TrimSpace(stringValue)
	stringValue = strings.ToLower(stringValue)

	re := regexp.MustCompile(marketplaceItemSlugRegexExpression)
	isValid := re.MatchString(stringValue)
	if !isValid {
		return "", errors.New("InvalidMarketplaceItemSlug")
	}

	return MarketplaceItemSlug(stringValue), nil
}

func NewMarketplaceItemSlugPanic(value interface{}) MarketplaceItemSlug {
	vo, err := NewMarketplaceItemSlug(value)
	if err != nil {
		panic(err)
	}

	return vo
}

func (vo MarketplaceItemSlug) String() string {
	return string(vo)
}
