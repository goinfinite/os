package entity

import (
	"github.com/speedianet/os/src/domain/valueObject"
)

type MarketplaceInstalledItem struct {
	Id               valueObject.MarketplaceInstalledItemId   `json:"id"`
	Name             valueObject.MarketplaceItemName          `json:"name"`
	Hostname         valueObject.Fqdn                         `json:"hostname"`
	Type             valueObject.MarketplaceItemType          `json:"type"`
	InstallDirectory valueObject.UnixFilePath                 `json:"installDirectory"`
	InstallUuid      valueObject.MarketplaceInstalledItemUuid `json:"installUuid"`
	ServiceNames     []valueObject.ServiceName                `json:"serviceNames"`
	Mappings         []Mapping                                `json:"mappings"`
	AvatarUrl        valueObject.Url                          `json:"avatarUrl"`
	CreatedAt        valueObject.UnixTime                     `json:"createdAt"`
	UpdatedAt        valueObject.UnixTime                     `json:"updatedAt"`
}

func NewMarketplaceInstalledItem(
	id valueObject.MarketplaceInstalledItemId,
	itemName valueObject.MarketplaceItemName,
	hostname valueObject.Fqdn,
	itemType valueObject.MarketplaceItemType,
	installDirectory valueObject.UnixFilePath,
	installUuid valueObject.MarketplaceInstalledItemUuid,
	serviceNames []valueObject.ServiceName,
	mappings []Mapping,
	avatarUrl valueObject.Url,
	createdAt valueObject.UnixTime,
	updatedAt valueObject.UnixTime,
) MarketplaceInstalledItem {
	return MarketplaceInstalledItem{
		Id:               id,
		Name:             itemName,
		Hostname:         hostname,
		Type:             itemType,
		InstallDirectory: installDirectory,
		InstallUuid:      installUuid,
		ServiceNames:     serviceNames,
		Mappings:         mappings,
		AvatarUrl:        avatarUrl,
		CreatedAt:        createdAt,
		UpdatedAt:        createdAt,
	}
}
