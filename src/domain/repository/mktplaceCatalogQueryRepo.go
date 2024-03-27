package repository

import "github.com/speedianet/os/src/domain/entity"

type MktplaceCatalogQueryRepo interface {
	GetItems() ([]entity.MarketplaceCatalogItem, error)
}
