package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/repository"
)

func ReadInstallableServices(
	servicesQueryRepo repository.ServicesQueryRepo,
) ([]entity.InstallableService, error) {
	installableServices, err := servicesQueryRepo.ReadInstallables()
	if err != nil {
		log.Printf("ReadInstallableServicesError: %s", err.Error())
		return installableServices, errors.New("ReadInstallableServicesInfraError")
	}

	return installableServices, nil
}
