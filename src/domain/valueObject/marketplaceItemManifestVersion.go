package valueObject

import (
	"errors"
	"slices"
	"strings"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

type MarketplaceItemManifestVersion string

var validMarketplaceItemManifestVersions = []string{
	"v1",
}

func NewMarketplaceItemManifestVersion(value interface{}) (
	version MarketplaceItemManifestVersion, err error,
) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return version, errors.New("MarketplaceItemManifestVersionMustBeString")
	}
	stringValue = strings.ToLower(stringValue)

	if !slices.Contains(validMarketplaceItemManifestVersions, stringValue) {
		return version, errors.New("InvalidMarketplaceItemManifestVersion")
	}

	return MarketplaceItemManifestVersion(stringValue), nil
}

func (vo MarketplaceItemManifestVersion) String() string {
	return string(vo)
}
