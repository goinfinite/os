package entity

import (
	"github.com/speedianet/os/src/domain/valueObject"
)

type MarketplaceCatalogItem struct {
	Id                 valueObject.MarketplaceCatalogItemId          `json:"id"`
	Name               valueObject.MarketplaceItemName               `json:"name"`
	Type               valueObject.MarketplaceItemType               `json:"type"`
	Description        valueObject.MarketplaceItemDescription        `json:"description"`
	ServiceNames       []valueObject.ServiceName                     `json:"services"`
	Mappings           []valueObject.MarketplaceItemMapping          `json:"mappings"`
	DataFields         []valueObject.MarketplaceCatalogItemDataField `json:"dataFields"`
	CmdSteps           []valueObject.MarketplaceItemCmdStep          `json:"-"`
	EstimatedSizeBytes valueObject.Byte                              `json:"estimatedSizeBytes"`
	AvatarUrl          valueObject.Url                               `json:"avatarUrl"`
	ScreenshotUrls     []valueObject.Url                             `json:"screenshotUrls"`
}

func NewMarketplaceCatalogItem(
	id valueObject.MarketplaceCatalogItemId,
	itemName valueObject.MarketplaceItemName,
	itemType valueObject.MarketplaceItemType,
	description valueObject.MarketplaceItemDescription,
	serviceNames []valueObject.ServiceName,
	mappings []valueObject.MarketplaceItemMapping,
	dataFields []valueObject.MarketplaceCatalogItemDataField,
	cmdSteps []valueObject.MarketplaceItemCmdStep,
	estimatedSizeBytes valueObject.Byte,
	avatarUrl valueObject.Url,
	screenshotUrls []valueObject.Url,
) MarketplaceCatalogItem {
	return MarketplaceCatalogItem{
		Id:                 id,
		Name:               itemName,
		Type:               itemType,
		Description:        description,
		ServiceNames:       serviceNames,
		Mappings:           mappings,
		DataFields:         dataFields,
		CmdSteps:           cmdSteps,
		EstimatedSizeBytes: estimatedSizeBytes,
		AvatarUrl:          avatarUrl,
		ScreenshotUrls:     screenshotUrls,
	}
}
