package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
	tkDto "github.com/goinfinite/tk/src/domain/dto"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

var (
	servicesDefaultSortBy        tkValueObject.PaginationSortBy        = "name"
	servicesDefaultSortDirection tkValueObject.PaginationSortDirection = "asc"

	ServicesDefaultPagination tkDto.Pagination = tkDto.Pagination{
		PageNumber:    0,
		ItemsPerPage:  10,
		SortBy:        &servicesDefaultSortBy,
		SortDirection: &servicesDefaultSortDirection,
	}
)

func ReadInstallableServices(
	servicesQueryRepo repository.ServicesQueryRepo,
	requestDto dto.ReadInstallableServicesItemsRequest,
) (responseDto dto.ReadInstallableServicesItemsResponse, err error) {
	responseDto, err = servicesQueryRepo.ReadInstallableItems(requestDto)
	if err != nil {
		slog.Error("ReadInstallableServicesError", slog.String("err", err.Error()))
		return responseDto, errors.New("ReadInstallableServicesInfraError")
	}

	return responseDto, nil
}
