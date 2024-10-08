package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/repository"
)

func ReadInstallableServices(
	servicesQueryRepo repository.ServicesQueryRepo,
) ([]entity.InstallableService, error) {
	installableServices, err := servicesQueryRepo.ReadInstallables()
	if err != nil {
		slog.Error("ReadInstallableServicesError", slog.Any("err", err))
		return installableServices, errors.New("ReadInstallableServicesInfraError")
	}

	return installableServices, nil
}
