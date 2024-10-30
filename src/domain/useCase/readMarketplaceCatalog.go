package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

var MarketplaceDefaultPagination dto.Pagination = dto.Pagination{
	PageNumber:   0,
	ItemsPerPage: 10,
}

func ReadMarketplaceCatalog(
	marketplaceQueryRepo repository.MarketplaceQueryRepo,
	readDto dto.ReadMarketplaceCatalogItemsRequest,
) (dto.ReadMarketplaceCatalogItemsResponse, error) {
	responseDto, err := marketplaceQueryRepo.ReadCatalogItems(readDto)
	if err != nil {
		slog.Error("ReadMarketplaceCatalogItemsError", slog.Any("error", err))
		return responseDto, errors.New("ReadMarketplaceCatalogItemsInfraError")
	}

	return responseDto, nil
}
