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
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	deleteDto dto.DeleteMarketplaceInstalledItem,
) error {
	_, err := marketplaceQueryRepo.ReadInstalledItemById(deleteDto.InstalledId)
	if err != nil {
		return errors.New("MarketplaceInstalledItemNotFound")
	}

	err = marketplaceCmdRepo.UninstallItem(deleteDto)
	if err != nil {
		slog.Error("UninstallMarketplaceItemError", slog.Any("error", err))
		return errors.New("UninstallMarketplaceItemInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		DeleteMarketplaceInstalledItem(deleteDto)

	return nil
}
