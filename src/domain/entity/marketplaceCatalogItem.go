package entity

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type MarketplaceCatalogItem struct {
	ManifestVersion      valueObject.MarketplaceItemManifestVersion    `json:"manifestVersion"`
	Id                   valueObject.MarketplaceItemId                 `json:"id"`
	Slugs                []valueObject.MarketplaceItemSlug             `json:"slugs"`
	Name                 valueObject.MarketplaceItemName               `json:"name"`
	Type                 valueObject.MarketplaceItemType               `json:"type"`
	Description          valueObject.MarketplaceItemDescription        `json:"description"`
	Services             []valueObject.ServiceNameWithVersion          `json:"services"`
	Mappings             []valueObject.MarketplaceItemMapping          `json:"mappings"`
	DataFields           []valueObject.MarketplaceCatalogItemDataField `json:"dataFields"`
	InstallTimeoutSecs   tkValueObject.UnixTime                        `json:"-"`
	InstallCmdSteps      []tkValueObject.UnixCommand                   `json:"-"`
	UninstallTimeoutSecs tkValueObject.UnixTime                        `json:"-"`
	UninstallCmdSteps    []tkValueObject.UnixCommand                   `json:"-"`
	UninstallFileNames   []tkValueObject.UnixFileName                  `json:"-"`
	EstimatedSizeBytes   tkValueObject.Byte                            `json:"estimatedSizeBytes"`
	AvatarUrl            tkValueObject.Url                             `json:"avatarUrl"`
	ScreenshotUrls       []tkValueObject.Url                           `json:"screenshotUrls"`
}

func NewMarketplaceCatalogItem(
	manifestVersion valueObject.MarketplaceItemManifestVersion,
	id valueObject.MarketplaceItemId,
	slugs []valueObject.MarketplaceItemSlug,
	itemName valueObject.MarketplaceItemName,
	itemType valueObject.MarketplaceItemType,
	description valueObject.MarketplaceItemDescription,
	services []valueObject.ServiceNameWithVersion,
	mappings []valueObject.MarketplaceItemMapping,
	dataFields []valueObject.MarketplaceCatalogItemDataField,
	installTimeoutSecs tkValueObject.UnixTime,
	installCmdSteps []tkValueObject.UnixCommand,
	uninstallTimeoutSecs tkValueObject.UnixTime,
	uninstallCmdSteps []tkValueObject.UnixCommand,
	uninstallFileNames []tkValueObject.UnixFileName,
	estimatedSizeBytes tkValueObject.Byte,
	avatarUrl tkValueObject.Url,
	screenshotUrls []tkValueObject.Url,
) MarketplaceCatalogItem {
	return MarketplaceCatalogItem{
		ManifestVersion:      manifestVersion,
		Id:                   id,
		Slugs:                slugs,
		Name:                 itemName,
		Type:                 itemType,
		Description:          description,
		Services:             services,
		Mappings:             mappings,
		DataFields:           dataFields,
		InstallTimeoutSecs:   installTimeoutSecs,
		InstallCmdSteps:      installCmdSteps,
		UninstallTimeoutSecs: uninstallTimeoutSecs,
		UninstallCmdSteps:    uninstallCmdSteps,
		UninstallFileNames:   uninstallFileNames,
		EstimatedSizeBytes:   estimatedSizeBytes,
		AvatarUrl:            avatarUrl,
		ScreenshotUrls:       screenshotUrls,
	}
}
