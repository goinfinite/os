package repository

import "github.com/speedianet/os/src/domain/entity"

type MarketplaceQueryRepo interface {
	Get() ([]entity.MarketplaceCatalogItem, error)
}
