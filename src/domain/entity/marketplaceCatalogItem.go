package entity

import (
	"github.com/goinfinite/os/src/domain/valueObject"
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
	InstallTimeoutSecs   valueObject.UnixTime                          `json:"-"`
	InstallCmdSteps      []valueObject.UnixCommand                     `json:"-"`
	UninstallTimeoutSecs valueObject.UnixTime                          `json:"-"`
	UninstallCmdSteps    []valueObject.UnixCommand                     `json:"-"`
	UninstallFileNames   []valueObject.UnixFileName                    `json:"-"`
	EstimatedSizeBytes   valueObject.Byte                              `json:"estimatedSizeBytes"`
	AvatarUrl            valueObject.Url                               `json:"avatarUrl"`
	ScreenshotUrls       []valueObject.Url                             `json:"screenshotUrls"`
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
	installTimeoutSecs valueObject.UnixTime,
	installCmdSteps []valueObject.UnixCommand,
	uninstallTimeoutSecs valueObject.UnixTime,
	uninstallCmdSteps []valueObject.UnixCommand,
	uninstallFileNames []valueObject.UnixFileName,
	estimatedSizeBytes valueObject.Byte,
	avatarUrl valueObject.Url,
	screenshotUrls []valueObject.Url,
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
