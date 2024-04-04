package entity

import (
	"github.com/speedianet/os/src/domain/valueObject"
)

type MarketplaceInstalledItem struct {
	Id            valueObject.MarketplaceItemId   `json:"id"`
	Name          valueObject.MarketplaceItemName `json:"name"`
	Type          valueObject.MarketplaceItemType `json:"type"`
	RootDirectory valueObject.UnixFilePath        `json:"rootDirectory"`
	Services      []valueObject.ServiceName       `json:"services"`
	Mappings      []Mapping                       `json:"mappings"`
	AvatarUrl     valueObject.Url                 `json:"avatarUrl"`
	CreatedAt     valueObject.UnixTime            `json:"createdAt"`
	UpdatedAt     valueObject.UnixTime            `json:"updatedAt"`
}

func NewMarketplaceInstalledItem(
	id valueObject.MarketplaceItemId,
	itemName valueObject.MarketplaceItemName,
	itemType valueObject.MarketplaceItemType,
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
