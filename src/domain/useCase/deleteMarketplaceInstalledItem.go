package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
)

func DeleteMarketplaceInstalledItem(
	marketplaceQueryRepo repository.MarketplaceQueryRepo,
	marketplaceCmdRepo repository.MarketplaceCmdRepo,
	dto dto.DeleteMarketplaceInstalledItem,
) error {
	_, err := marketplaceQueryRepo.GetInstalledItemById(dto.InstalledId)
	if err != nil {
		return errors.New("MarketplaceInstalledItemNotFound")
	}

	err = marketplaceCmdRepo.UninstallItem(dto)
	if err != nil {
		log.Printf("UninstallMarketplaceItemError: %s", err.Error())
		return errors.New("UninstallMarketplaceItemInfraError")
	}

	return nil
}
