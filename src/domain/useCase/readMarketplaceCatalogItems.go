package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
)

var MarketplaceDefaultSortBy valueObject.PaginationSortBy = valueObject.PaginationSortBy("id")

var MarketplaceDefaultPagination dto.Pagination = dto.Pagination{
	PageNumber:   0,
	ItemsPerPage: 50,
	SortBy:       &MarketplaceDefaultSortBy,
}

func ReadMarketplaceCatalogItems(
	marketplaceQueryRepo repository.MarketplaceQueryRepo,
	requestDto dto.ReadMarketplaceCatalogItemsRequest,
) (responseDto dto.ReadMarketplaceCatalogItemsResponse, err error) {
	responseDto, err = marketplaceQueryRepo.ReadCatalogItems(requestDto)
	if err != nil {
		slog.Error("ReadMarketplaceCatalogItemsError", slog.String("err", err.Error()))
		return responseDto, errors.New("ReadMarketplaceCatalogItemsInfraError")
	}

	return responseDto, nil
}
