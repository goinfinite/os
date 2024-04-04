package repository

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type MarketplaceQueryRepo interface {
	GetItems() ([]entity.MarketplaceCatalogItem, error)
	GetItemById(id valueObject.MarketplaceItemId) (entity.MarketplaceCatalogItem, error)
	GetInstalledItems() ([]entity.MarketplaceInstalledItem, error)
}
