package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

func ReadInstalledServices(
	servicesQueryRepo repository.ServicesQueryRepo,
	requestDto dto.ReadInstalledServicesItemsRequest,
) (responseDto dto.ReadInstalledServicesItemsResponse, err error) {
	responseDto, err = servicesQueryRepo.ReadInstalledItems(requestDto)
	if err != nil {
		slog.Error("ReadInstalledServicesError", slog.Any("error", err))
		return responseDto, errors.New("ReadInstalledServicesInfraError")
	}

	return responseDto, nil
}
