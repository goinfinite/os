package dto

import (
	"github.com/speedianet/os/src/domain/valueObject"
)

type InstallMarketplaceCatalogItem struct {
	Id         valueObject.MarketplaceCatalogItemId
	Hostname   valueObject.Fqdn
	UrlPath    *valueObject.UrlPath
	DataFields []valueObject.MarketplaceInstallableItemDataField
}

func NewInstallMarketplaceCatalogItem(
	id valueObject.MarketplaceCatalogItemId,
	hostname valueObject.Fqdn,
	urlPath *valueObject.UrlPath,
	dataFields []valueObject.MarketplaceInstallableItemDataField,
) InstallMarketplaceCatalogItem {
	return InstallMarketplaceCatalogItem{
		Id:         id,
		Hostname:   hostname,
		UrlPath:    urlPath,
		DataFields: dataFields,
	}
}
