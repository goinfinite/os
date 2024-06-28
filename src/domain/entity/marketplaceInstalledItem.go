package entity

import (
	"github.com/speedianet/os/src/domain/valueObject"
)

type MarketplaceInstalledItem struct {
	Id               valueObject.MarketplaceItemId            `json:"id"`
	Name             valueObject.MarketplaceItemName          `json:"name"`
	Hostname         valueObject.Fqdn                         `json:"hostname"`
	Type             valueObject.MarketplaceItemType          `json:"type"`
	UrlPath          valueObject.UrlPath                      `json:"urlPath"`
	InstallDirectory valueObject.UnixFilePath                 `json:"installDirectory"`
	InstallUuid      valueObject.MarketplaceInstalledItemUuid `json:"installUuid"`
	Services         []valueObject.ServiceNameWithVersion     `json:"services"`
	Mappings         []Mapping                                `json:"mappings"`
	AvatarUrl        valueObject.Url                          `json:"avatarUrl"`
	AppSlug          valueObject.MarketplaceItemSlug          `json:"-"`
	CreatedAt        valueObject.UnixTime                     `json:"createdAt"`
	UpdatedAt        valueObject.UnixTime                     `json:"updatedAt"`
}

func NewMarketplaceInstalledItem(
	id valueObject.MarketplaceItemId,
	itemName valueObject.MarketplaceItemName,
	hostname valueObject.Fqdn,
	itemType valueObject.MarketplaceItemType,
	urlPath valueObject.UrlPath,
	installDirectory valueObject.UnixFilePath,
	installUuid valueObject.MarketplaceInstalledItemUuid,
	services []valueObject.ServiceNameWithVersion,
	mappings []Mapping,
	avatarUrl valueObject.Url,
	appSlug valueObject.MarketplaceItemSlug,
	createdAt valueObject.UnixTime,
	updatedAt valueObject.UnixTime,
) MarketplaceInstalledItem {
	return MarketplaceInstalledItem{
		Id:               id,
		Name:             itemName,
		Hostname:         hostname,
		Type:             itemType,
		UrlPath:          urlPath,
		InstallDirectory: installDirectory,
		InstallUuid:      installUuid,
		Services:         services,
		Mappings:         mappings,
		AvatarUrl:        avatarUrl,
		AppSlug:          appSlug,
		CreatedAt:        createdAt,
		UpdatedAt:        createdAt,
	}
}
