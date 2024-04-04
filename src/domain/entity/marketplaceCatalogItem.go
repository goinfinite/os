package entity

import (
	"github.com/speedianet/os/src/domain/valueObject"
)

type MarketplaceCatalogItem struct {
	Id                 valueObject.MarketplaceItemId            `json:"id" yaml:"id"`
	Name               valueObject.MarketplaceItemName          `json:"name" yaml:"name"`
	Type               valueObject.MarketplaceItemType          `json:"type" yaml:"type"`
	Description        valueObject.MarketplaceItemDescription   `json:"description" yaml:"description"`
	Services           []valueObject.ServiceName                `json:"services" yaml:"services"`
	Mappings           []MarketplaceMapping                     `json:"mappings" yaml:"mappings"`
	DataFields         []valueObject.DataField                  `json:"dataFields" yaml:"dataFields"`
	CmdSteps           []valueObject.MarketplaceItemInstallStep `json:"cmdSteps" yaml:"cmdSteps"`
	EstimatedSizeBytes valueObject.Byte                         `json:"estimatedSizeBytes" yaml:"estimatedSizeBytes"`
	AvatarUrl          valueObject.Url                          `json:"avatarUrl" yaml:"avatarUrl"`
	ScreenshotUrls     []valueObject.Url                        `json:"screenshotUrls" yaml:"screenshotUrls"`
}
