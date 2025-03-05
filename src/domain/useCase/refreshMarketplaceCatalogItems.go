package useCase

import (
	"log/slog"

	"github.com/goinfinite/os/src/domain/repository"
)

const RefreshMarketplaceCatalogItemsAmountPerDay int = 1

func RefreshMarketplaceCatalogItems(marketplaceCmdRepo repository.MarketplaceCmdRepo) {
	err := marketplaceCmdRepo.RefreshCatalogItems()
	if err != nil {
		slog.Error("RefreshMarketplaceCatalogItemsError", slog.String("err", err.Error()))
	}
}
