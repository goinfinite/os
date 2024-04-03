package repository

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type MktplaceQueryRepo interface {
	GetItems() ([]entity.MarketplaceCatalogItem, error)
	GetItemById(id valueObject.MktplaceItemId) (entity.MarketplaceCatalogItem, error)
	GetInstalledItems() ([]entity.MarketplaceInstalledItem, error)
}
