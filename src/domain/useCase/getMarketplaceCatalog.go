package useCase

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/repository"
)

func GetMarketplaceCatalog(
	mktplaceCatalogQueryRepo repository.MktplaceCatalogQueryRepo,
) ([]entity.MarketplaceCatalogItem, error) {
	return mktplaceCatalogQueryRepo.GetItems()
}
