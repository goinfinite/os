package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/repository"
)

func ReadMarketplaceCatalog(
	marketplaceQueryRepo repository.MarketplaceQueryRepo,
) ([]entity.MarketplaceCatalogItem, error) {
	catalogItems, err := marketplaceQueryRepo.ReadCatalogItems()
	if err != nil {
		slog.Error("ReadMarketplaceCatalogItemsError", slog.Any("err", err))
		return nil, errors.New("ReadMarketplaceCatalogItemsInfraError")
	}

	return catalogItems, nil
}
