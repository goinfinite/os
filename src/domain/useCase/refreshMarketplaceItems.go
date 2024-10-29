package useCase

import (
	"log/slog"

	"github.com/goinfinite/os/src/domain/repository"
)

func RefreshMarketplaceItems(marketplaceCmdRepo repository.MarketplaceCmdRepo) {
	err := marketplaceCmdRepo.RefreshItems()
	if err != nil {
		slog.Error("RefreshMarketplaceItemsInfraError", slog.Any("error", err))
	}
}
