package repository

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type MarketplaceQueryRepo interface {
	GetCatalogItems() ([]entity.MarketplaceCatalogItem, error)
	GetCatalogItemById(id valueObject.MarketplaceItemId) (entity.MarketplaceCatalogItem, error)
	GetInstalledItems() ([]entity.MarketplaceInstalledItem, error)
}
