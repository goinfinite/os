package useCase

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/repository"
)

func GetMarketplaceCatalog(
	mktplaceQueryRepo repository.MktplaceQueryRepo,
) ([]entity.MarketplaceCatalogItem, error) {
	return mktplaceQueryRepo.GetItems()
}
