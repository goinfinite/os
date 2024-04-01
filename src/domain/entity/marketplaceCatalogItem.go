package entity

import (
	"github.com/speedianet/os/src/domain/valueObject"
)

type MarketplaceCatalogItem struct {
	Id                 valueObject.MktplaceItemId            `json:"id"`
	Name               valueObject.MktplaceItemName          `json:"name"`
	Type               valueObject.MktplaceItemType          `json:"type"`
	Description        valueObject.MktplaceItemDescription   `json:"description"`
	Services           []valueObject.ServiceName             `json:"services"`
	Mappings           []MarketplaceMapping                  `json:"mappings"`
	DataFields         []valueObject.DataField               `json:"-"`
	CmdSteps           []valueObject.MktplaceItemInstallStep `json:"cmdSteps"`
	EstimatedSizeBytes valueObject.Byte                      `json:"estimatedSizeBytes"`
	AvatarUrl          valueObject.Url                       `json:"avatarUrl"`
	ScreenshotUrls     []valueObject.Url                     `json:"screenshotUrls"`
}
