package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/repository"
)

func ReadMarketplaceInstalledItems(
	marketplaceQueryRepo repository.MarketplaceQueryRepo,
) ([]entity.MarketplaceInstalledItem, error) {
	installedItems, err := marketplaceQueryRepo.ReadInstalledItems()
	if err != nil {
		slog.Error("ReadMarketplaceInstalledItemsError", slog.Any("error", err))
		return nil, errors.New("ReadMarketplaceInstalledItemsInfraError")
	}

	return installedItems, nil
}
