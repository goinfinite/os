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
	Mappings      []Mapping                    `json:"mappings"`
	AvatarUrl     valueObject.Url              `json:"avatarUrl"`
	CreatedAt     valueObject.UnixTime         `json:"createdAt"`
	UpdatedAt     valueObject.UnixTime         `json:"updatedAt"`
}

func NewMarketplaceInstalledItem(
	id valueObject.MktplaceItemId,
	itemName valueObject.MktplaceItemName,
	itemType valueObject.MktplaceItemType,
	rootDirectory valueObject.UnixFilePath,
	services []valueObject.ServiceName,
	mappings []Mapping,
	avatarUrl valueObject.Url,
	createdAt valueObject.UnixTime,
	updatedAt valueObject.UnixTime,
) MarketplaceInstalledItem {
	return MarketplaceInstalledItem{
		Id:            id,
		Name:          itemName,
		Type:          itemType,
		RootDirectory: rootDirectory,
		Services:      services,
		Mappings:      mappings,
		AvatarUrl:     avatarUrl,
		CreatedAt:     createdAt,
		UpdatedAt:     createdAt,
	}
}
