package useCase

import (
	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
)

func DeleteMarketplaceInstalledItem(
	marketplaceQueryRepo repository.MarketplaceQueryRepo,
	marketplaceCmdRepo repository.MarketplaceCmdRepo,
	installedId valueObject.MarketplaceInstalledItemId,
) error {
	return nil
}
