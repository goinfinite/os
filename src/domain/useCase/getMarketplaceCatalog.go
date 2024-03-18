package useCase

import (
	"github.com/speedianet/os/src/domain/repository"
)

func GetMarketplaceCatalog(
	mktplaceQueryRepo repository.MarketplaceQueryRepo,
) error {
	return mktplaceQueryRepo.Get()
}
