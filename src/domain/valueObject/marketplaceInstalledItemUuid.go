package valueObject

import (
	"errors"
	"regexp"
	"strings"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

type MarketplaceInstalledItemUuid string

var marketplaceInstalledItemUuidRegex = regexp.MustCompile(`^\w{10,16}$`)

func NewMarketplaceInstalledItemUuid(value any) (
	marketplaceInstalledItemUuid MarketplaceInstalledItemUuid, err error,
) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
	if err != nil {
		return marketplaceInstalledItemUuid, errors.New(
			"MarketplaceInstalledItemUuidMustBeString",
		)
	}
	stringValue = strings.ToLower(stringValue)

	if !marketplaceInstalledItemUuidRegex.MatchString(stringValue) {
		return marketplaceInstalledItemUuid, errors.New(
			"InvalidMarketplaceInstalledItemUuid",
		)
	}

	return MarketplaceInstalledItemUuid(stringValue), nil
}

func (vo MarketplaceInstalledItemUuid) String() string {
	return string(vo)
}
