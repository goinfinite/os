package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

func DeleteMarketplaceInstalledItem(
	marketplaceQueryRepo repository.MarketplaceQueryRepo,
	marketplaceCmdRepo repository.MarketplaceCmdRepo,
	deleteDto dto.DeleteMarketplaceInstalledItem,
) error {
	readInstalledItemDto := dto.ReadMarketplaceInstalledItemsRequest{
		MarketplaceInstalledItemId: &deleteDto.InstalledId,
	}
	_, err := marketplaceQueryRepo.ReadOneInstalledItem(readInstalledItemDto)
	if err != nil {
		return errors.New("MarketplaceInstalledItemNotFound")
	}

	err = marketplaceCmdRepo.UninstallItem(deleteDto)
	if err != nil {
		slog.Error("UninstallMarketplaceItemError", slog.Any("error", err))
		return errors.New("UninstallMarketplaceItemInfraError")
	}

	return nil
}
