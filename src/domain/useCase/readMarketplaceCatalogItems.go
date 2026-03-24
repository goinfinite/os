package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
	tkDto "github.com/goinfinite/tk/src/domain/dto"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

var MarketplaceDefaultSortBy tkValueObject.PaginationSortBy = tkValueObject.PaginationSortBy("id")

var MarketplaceDefaultPagination tkDto.Pagination = tkDto.Pagination{
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
