package entity

import (
	"github.com/speedianet/os/src/domain/valueObject"
)

type MarketplaceInstalledItem struct {
	Id               valueObject.MarketplaceItemId   `json:"id"`
	Name             valueObject.MarketplaceItemName `json:"name"`
	Type             valueObject.MarketplaceItemType `json:"type"`
	InstallDirectory valueObject.UnixFilePath        `json:"installDirectory"`
	ServiceNames     []valueObject.ServiceName       `json:"serviceNames"`
	Mappings         []Mapping                       `json:"mappings"`
	AvatarUrl        valueObject.Url                 `json:"avatarUrl"`
	CreatedAt        valueObject.UnixTime            `json:"createdAt"`
	UpdatedAt        valueObject.UnixTime            `json:"updatedAt"`
}

func NewMarketplaceInstalledItem(
	id valueObject.MarketplaceItemId,
	itemName valueObject.MarketplaceItemName,
	itemType valueObject.MarketplaceItemType,
	installDirectory valueObject.UnixFilePath,
	serviceNames []valueObject.ServiceName,
	mappings []Mapping,
	avatarUrl valueObject.Url,
	createdAt valueObject.UnixTime,
	updatedAt valueObject.UnixTime,
) MarketplaceInstalledItem {
	return MarketplaceInstalledItem{
		Id:               id,
		Name:             itemName,
		Type:             itemType,
		InstallDirectory: installDirectory,
		ServiceNames:     serviceNames,
		Mappings:         mappings,
		AvatarUrl:        avatarUrl,
		CreatedAt:        createdAt,
		UpdatedAt:        createdAt,
	}
}
