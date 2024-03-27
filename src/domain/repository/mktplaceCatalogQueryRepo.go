package repository

import "github.com/speedianet/os/src/domain/entity"

type MktplaceCatalogQueryRepo interface {
	Get() ([]entity.MarketplaceCatalogItem, error)
}
