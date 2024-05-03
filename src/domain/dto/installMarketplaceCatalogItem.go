package dto

import (
	"github.com/speedianet/os/src/domain/valueObject"
)

type InstallMarketplaceCatalogItem struct {
	Id           valueObject.MarketplaceCatalogItemId
	Hostname     valueObject.Fqdn
	UrlDirectory *valueObject.UnixFilePath
	DataFields   []valueObject.MarketplaceInstallableItemDataField
}

func NewInstallMarketplaceCatalogItem(
	id valueObject.MarketplaceCatalogItemId,
	hostname valueObject.Fqdn,
	urlDirectory *valueObject.UnixFilePath,
	dataFields []valueObject.MarketplaceInstallableItemDataField,
) InstallMarketplaceCatalogItem {
	return InstallMarketplaceCatalogItem{
		Id:           id,
		Hostname:     hostname,
		UrlDirectory: urlDirectory,
		DataFields:   dataFields,
	}
}
