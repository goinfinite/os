package dto

import (
	"github.com/speedianet/os/src/domain/valueObject"
)

type InstallMarketplaceCatalogItem struct {
	Id            valueObject.MktplaceItemId
	Hostname      valueObject.Fqdn
	RootDirectory valueObject.UnixFilePath
	DataFields    []valueObject.DataField
}

func NewInstallMarketplaceCatalogItem(
	id valueObject.MktplaceItemId,
	hostname valueObject.Fqdn,
	rootDirectory valueObject.UnixFilePath,
	dataFields []valueObject.DataField,
) InstallMarketplaceCatalogItem {
	return InstallMarketplaceCatalogItem{
		Id:            id,
		Hostname:      hostname,
		RootDirectory: rootDirectory,
		DataFields:    dataFields,
	}
}
