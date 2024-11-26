package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

var ServicesDefaultPagination dto.Pagination = dto.Pagination{
	PageNumber:   0,
	ItemsPerPage: 10,
}

func ReadInstallableServices(
	servicesQueryRepo repository.ServicesQueryRepo,
	requestDto dto.ReadInstallableServicesItemsRequest,
) (responseDto dto.ReadInstallableServicesItemsResponse, err error) {
	responseDto, err = servicesQueryRepo.ReadInstallableItems(requestDto)
	if err != nil {
		slog.Error("ReadInstallableServicesError", slog.Any("error", err))
		return responseDto, errors.New("ReadInstallableServicesInfraError")
	}

	return responseDto, nil
}
