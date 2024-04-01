package dto

import (
	"github.com/speedianet/os/src/domain/valueObject"
)

type InstallMarketplaceItem struct {
	Name          valueObject.MktplaceItemName
	Type          valueObject.MktplaceItemType
	RootDirectory valueObject.UnixFilePath
	DataFields    []valueObject.DataField
}

func NewInstallMarketplaceItem(
	itemName valueObject.MktplaceItemName,
	itemType valueObject.MktplaceItemType,
	rootDirectory valueObject.UnixFilePath,
	dataFields []valueObject.DataField,
) InstallMarketplaceItem {
	return InstallMarketplaceItem{
		Name:          itemName,
		Type:          itemType,
		RootDirectory: rootDirectory,
		DataFields:    dataFields,
	}
}
