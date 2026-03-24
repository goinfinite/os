package entity

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type MarketplaceInstalledItem struct {
	Id               valueObject.MarketplaceItemId            `json:"id"`
	Name             valueObject.MarketplaceItemName          `json:"name"`
	Hostname         tkValueObject.Fqdn                      `json:"hostname"`
	Type             valueObject.MarketplaceItemType          `json:"type"`
	UrlPath          valueObject.UrlPath                      `json:"urlPath"`
	InstallDirectory tkValueObject.UnixAbsoluteFilePath       `json:"installDirectory"`
	InstallUuid      valueObject.MarketplaceInstalledItemUuid `json:"installUuid"`
	Services         []valueObject.ServiceNameWithVersion     `json:"services"`
	Mappings         []Mapping                                `json:"mappings"`
	AvatarUrl        tkValueObject.Url                        `json:"avatarUrl"`
	Slug             valueObject.MarketplaceItemSlug          `json:"-"`
	CreatedAt        tkValueObject.UnixTime                   `json:"createdAt"`
	UpdatedAt        tkValueObject.UnixTime                   `json:"updatedAt"`
}

func NewMarketplaceInstalledItem(
	id valueObject.MarketplaceItemId,
	itemName valueObject.MarketplaceItemName,
	hostname tkValueObject.Fqdn,
	itemType valueObject.MarketplaceItemType,
	urlPath valueObject.UrlPath,
	installDirectory tkValueObject.UnixAbsoluteFilePath,
	installUuid valueObject.MarketplaceInstalledItemUuid,
	services []valueObject.ServiceNameWithVersion,
	mappings []Mapping,
	avatarUrl tkValueObject.Url,
	slug valueObject.MarketplaceItemSlug,
	createdAt tkValueObject.UnixTime,
	updatedAt tkValueObject.UnixTime,
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
		Slug:             slug,
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
	}
}
