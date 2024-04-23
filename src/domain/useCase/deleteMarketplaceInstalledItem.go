package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
)

func DeleteMarketplaceInstalledItem(
	marketplaceQueryRepo repository.MarketplaceQueryRepo,
	marketplaceCmdRepo repository.MarketplaceCmdRepo,
	installedId valueObject.MarketplaceInstalledItemId,
	shouldUninstallServices bool,
) error {
	_, err := marketplaceQueryRepo.GetInstalledItemById(installedId)
	if err != nil {
		return errors.New("MarketplaceInstalledItemNotFound")
	}

	err = marketplaceCmdRepo.UninstallItem(installedId, shouldUninstallServices)
	if err != nil {
		log.Printf("UninstallMarketplaceItemError: %s", err.Error())
		return errors.New("UninstallMarketplaceItemInfraError")
	}

	return nil
}
