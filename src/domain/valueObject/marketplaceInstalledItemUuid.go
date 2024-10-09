package valueObject

import (
	"errors"
	"regexp"
	"strings"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

type MarketplaceInstalledItemUuid string

const marketplaceInstalledItemUuidRegexExpression = `^\w{10,16}$`

func NewMarketplaceInstalledItemUuid(value interface{}) (
	marketplaceInstalledItemUuid MarketplaceInstalledItemUuid, err error,
) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return marketplaceInstalledItemUuid, errors.New(
			"MarketplaceInstalledItemUuidMustBeString",
		)
	}
	stringValue = strings.ToLower(stringValue)

	re := regexp.MustCompile(marketplaceInstalledItemUuidRegexExpression)
	if !re.MatchString(stringValue) {
		return marketplaceInstalledItemUuid, errors.New(
			"InvalidMarketplaceInstalledItemUuid",
		)
	}

	return MarketplaceInstalledItemUuid(stringValue), nil
}

func (vo MarketplaceInstalledItemUuid) String() string {
	return string(vo)
}
