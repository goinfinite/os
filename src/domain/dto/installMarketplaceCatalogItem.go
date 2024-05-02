package dto

import (
	"github.com/speedianet/os/src/domain/valueObject"
)

type InstallMarketplaceCatalogItem struct {
	Id         valueObject.MarketplaceCatalogItemId
	Hostname   valueObject.Fqdn
	Directory  *valueObject.UnixFilePath
	DataFields []valueObject.MarketplaceInstallableItemDataField
}

func NewInstallMarketplaceCatalogItem(
	id valueObject.MarketplaceCatalogItemId,
	hostname valueObject.Fqdn,
	installDirectory *valueObject.UnixFilePath,
	dataFields []valueObject.MarketplaceInstallableItemDataField,
) InstallMarketplaceCatalogItem {
	return InstallMarketplaceCatalogItem{
		Id:         id,
		Hostname:   hostname,
		Directory:  installDirectory,
		DataFields: dataFields,
	}
}
