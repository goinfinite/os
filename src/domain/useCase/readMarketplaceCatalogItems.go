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

func ReadMarketplaceCatalogItems(
	marketplaceQueryRepo repository.MarketplaceQueryRepo,
	requestDto dto.ReadMarketplaceCatalogItemsRequest,
) (responseDto dto.ReadMarketplaceCatalogItemsResponse, err error) {
	responseDto, err = marketplaceQueryRepo.ReadCatalogItems(requestDto)
	if err != nil {
		slog.Error("ReadMarketplaceCatalogItemsError", slog.Any("error", err))
		return responseDto, errors.New("ReadMarketplaceCatalogItemsInfraError")
	}

	return responseDto, nil
}
