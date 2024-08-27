package useCase

import (
	"errors"
	"log/slog"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/repository"
)

func ReadMarketplaceInstalledItems(
	marketplaceQueryRepo repository.MarketplaceQueryRepo,
) ([]entity.MarketplaceInstalledItem, error) {
	installedItems, err := marketplaceQueryRepo.ReadInstalledItems()
	if err != nil {
		slog.Error("ReadMarketplaceInstalledItemsError", slog.Any("err", err))
		return nil, errors.New("ReadMarketplaceInstalledItemsInfraError")
	}

	return installedItems, nil
}
