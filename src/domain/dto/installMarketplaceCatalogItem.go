package dto

import (
	"github.com/speedianet/os/src/domain/valueObject"
)

type InstallMarketplaceCatalogItem struct {
	Id               valueObject.MarketplaceItemId
	Hostname         valueObject.Fqdn
	InstallDirectory *valueObject.UnixFilePath
	DataFields       []valueObject.MarketplaceCatalogItemDataField
}

func NewInstallMarketplaceCatalogItem(
	id valueObject.MarketplaceItemId,
	hostname valueObject.Fqdn,
	installDirectory *valueObject.UnixFilePath,
	dataFields []valueObject.MarketplaceCatalogItemDataField,
) InstallMarketplaceCatalogItem {
	return InstallMarketplaceCatalogItem{
		Id:               id,
		Hostname:         hostname,
		InstallDirectory: installDirectory,
		DataFields:       dataFields,
	}
}
