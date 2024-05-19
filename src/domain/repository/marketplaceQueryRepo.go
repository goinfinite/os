package repository

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type MarketplaceQueryRepo interface {
	ReadCatalogItems() ([]entity.MarketplaceCatalogItem, error)
	ReadCatalogItemById(
		catalogId valueObject.MarketplaceItemId,
	) (entity.MarketplaceCatalogItem, error)
	ReadInstalledItems() ([]entity.MarketplaceInstalledItem, error)
	ReadInstalledItemById(
		installedId valueObject.MarketplaceItemId,
	) (entity.MarketplaceInstalledItem, error)
}
