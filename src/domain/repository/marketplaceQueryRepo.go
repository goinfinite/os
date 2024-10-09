package repository

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type MarketplaceQueryRepo interface {
	ReadCatalogItems() ([]entity.MarketplaceCatalogItem, error)
	ReadCatalogItemById(
		catalogId valueObject.MarketplaceItemId,
	) (entity.MarketplaceCatalogItem, error)
	ReadCatalogItemBySlug(
		slug valueObject.MarketplaceItemSlug,
	) (entity.MarketplaceCatalogItem, error)
	ReadInstalledItems() ([]entity.MarketplaceInstalledItem, error)
	ReadInstalledItemById(
		installedId valueObject.MarketplaceItemId,
	) (entity.MarketplaceInstalledItem, error)
}
