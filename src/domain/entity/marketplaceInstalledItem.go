package entity

import (
	"github.com/speedianet/os/src/domain/valueObject"
)

type MarketplaceInstalledItem struct {
	Id            valueObject.MktplaceItemId   `json:"id"`
	Name          valueObject.MktplaceItemName `json:"name"`
	Type          valueObject.MktplaceItemType `json:"type"`
	RootDirectory valueObject.UnixFilePath     `json:"rootDirectory"`
	Services      []valueObject.ServiceName    `json:"services"`
	MappingsIds   []valueObject.MappingId      `json:"mappingsIds"`
	AvatarUrl     valueObject.Url              `json:"avatarUrl"`
	CreatedAt     valueObject.UnixTime         `json:"createdAt"`
	UpdatedAt     valueObject.UnixTime         `json:"updatedAt"`
}
