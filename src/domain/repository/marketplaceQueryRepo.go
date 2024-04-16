package repository

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type MarketplaceQueryRepo interface {
	GetCatalogItems() ([]entity.MarketplaceCatalogItem, error)
	GetCatalogItemById(
		id valueObject.MarketplaceCatalogItemId,
	) (entity.MarketplaceCatalogItem, error)
	GetInstalledItems() ([]entity.MarketplaceInstalledItem, error)
	GetInstalledItemById(
		id valueObject.MarketplaceInstalledItemId,
	) (entity.MarketplaceInstalledItem, error)
}
