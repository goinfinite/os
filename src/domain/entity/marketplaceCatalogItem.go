package entity

import (
	"github.com/goinfinite/os/src/domain/valueObject"
)

type MarketplaceCatalogItem struct {
	Id                 valueObject.MarketplaceItemId                 `json:"id"`
	Slugs              []valueObject.MarketplaceItemSlug             `json:"slugs"`
	Name               valueObject.MarketplaceItemName               `json:"name"`
	Type               valueObject.MarketplaceItemType               `json:"type"`
	Description        valueObject.MarketplaceItemDescription        `json:"description"`
	Services           []valueObject.ServiceNameWithVersion          `json:"services"`
	Mappings           []valueObject.MarketplaceItemMapping          `json:"mappings"`
	DataFields         []valueObject.MarketplaceCatalogItemDataField `json:"dataFields"`
	InstallCmdSteps    []valueObject.UnixCommand                     `json:"-"`
	UninstallCmdSteps  []valueObject.UnixCommand                     `json:"-"`
	UninstallFileNames []valueObject.UnixFileName                    `json:"-"`
	EstimatedSizeBytes valueObject.Byte                              `json:"estimatedSizeBytes"`
	AvatarUrl          valueObject.Url                               `json:"avatarUrl"`
	ScreenshotUrls     []valueObject.Url                             `json:"screenshotUrls"`
}

func NewMarketplaceCatalogItem(
	id valueObject.MarketplaceItemId,
	slugs []valueObject.MarketplaceItemSlug,
	itemName valueObject.MarketplaceItemName,
	itemType valueObject.MarketplaceItemType,
	description valueObject.MarketplaceItemDescription,
	services []valueObject.ServiceNameWithVersion,
	mappings []valueObject.MarketplaceItemMapping,
	dataFields []valueObject.MarketplaceCatalogItemDataField,
	installCmdSteps []valueObject.UnixCommand,
	uninstallCmdSteps []valueObject.UnixCommand,
	uninstallFileNames []valueObject.UnixFileName,
	estimatedSizeBytes valueObject.Byte,
	avatarUrl valueObject.Url,
	screenshotUrls []valueObject.Url,
) MarketplaceCatalogItem {
	return MarketplaceCatalogItem{
		Id:                 id,
		Slugs:              slugs,
		Name:               itemName,
		Type:               itemType,
		Description:        description,
		Services:           services,
		Mappings:           mappings,
		DataFields:         dataFields,
		InstallCmdSteps:    installCmdSteps,
		UninstallCmdSteps:  uninstallCmdSteps,
		UninstallFileNames: uninstallFileNames,
		EstimatedSizeBytes: estimatedSizeBytes,
		AvatarUrl:          avatarUrl,
		ScreenshotUrls:     screenshotUrls,
	}
}
