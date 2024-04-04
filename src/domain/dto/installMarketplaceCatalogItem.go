package dto

import (
	"github.com/speedianet/os/src/domain/valueObject"
)

type InstallMarketplaceCatalogItem struct {
	Id            valueObject.MarketplaceItemId
	Hostname      valueObject.Fqdn
	RootDirectory valueObject.UnixFilePath
	DataFields    []valueObject.MarketplaceItemDataField
}

func NewInstallMarketplaceCatalogItem(
	id valueObject.MarketplaceItemId,
	hostname valueObject.Fqdn,
	rootDirectory valueObject.UnixFilePath,
	dataFields []valueObject.MarketplaceItemDataField,
) InstallMarketplaceCatalogItem {
	return InstallMarketplaceCatalogItem{
		Id:            id,
		Hostname:      hostname,
		RootDirectory: rootDirectory,
		DataFields:    dataFields,
	}
}
