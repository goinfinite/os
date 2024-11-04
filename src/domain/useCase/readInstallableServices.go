package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/repository"
)

var ServicesDefaultPagination dto.Pagination = dto.Pagination{
	PageNumber:   0,
	ItemsPerPage: 10,
}

func ReadInstallableServices(
	servicesQueryRepo repository.ServicesQueryRepo,
	readDto dto.ReadInstallableServicesItemsRequest,
) ([]entity.InstallableService, error) {
	installableServices, err := servicesQueryRepo.ReadInstallableItems()
	if err != nil {
		slog.Error("ReadInstallableServicesError", slog.Any("error", err))
		return installableServices, errors.New("ReadInstallableServicesInfraError")
	}

	return installableServices, nil
}
