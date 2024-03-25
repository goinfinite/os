package useCase

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/repository"
)

func GetMarketplaceCatalog(
	mktplaceQueryRepo repository.MarketplaceQueryRepo,
) ([]entity.MarketplaceCatalogItem, error) {
	return mktplaceQueryRepo.Get()
}
