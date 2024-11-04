package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

func ReadInstalledServices(
	servicesQueryRepo repository.ServicesQueryRepo,
	readDto dto.ReadInstalledServicesItemsRequest,
) (dto.ReadInstalledServicesItemsResponse, error) {
	installedServices, err := servicesQueryRepo.ReadInstalledItems(readDto)
	if err != nil {
		slog.Error("ReadInstalledServicesError", slog.Any("error", err))
		return installedServices, errors.New("ReadInstalledServicesInfraError")
	}

	return installedServices, nil
}
