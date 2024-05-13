package valueObject

import (
	"errors"
	"regexp"
	"strings"
)

type MarketplaceInstalledItemUuid string

const marketplaceInstalledItemUuidRegexExpression = `^\w{10,16}$`

func NewMarketplaceInstalledItemUuid(
	value string,
) (MarketplaceInstalledItemUuid, error) {
	value = strings.ToLower(value)

	min := MarketplaceInstalledItemUuid(value)
	if !min.isValid() {
		return "", errors.New("InvalidMarketplaceInstalledItemUuid")
	}

	return min, nil
}

func NewMarketplaceInstalledItemUuidPanic(
	value string,
) MarketplaceInstalledItemUuid {
	min, err := NewMarketplaceInstalledItemUuid(value)
	if err != nil {
		panic(err)
	}

	return min
}

func (min MarketplaceInstalledItemUuid) isValid() bool {
	marketplaceInstalledItemUuidCompiledRegex := regexp.MustCompile(
		marketplaceInstalledItemUuidRegexExpression,
	)
	return marketplaceInstalledItemUuidCompiledRegex.MatchString(string(min))
}

func (min MarketplaceInstalledItemUuid) String() string {
	return string(min)
}
