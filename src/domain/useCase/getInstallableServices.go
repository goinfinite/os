package useCase

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/repository"
)

func GetInstallableServices(
	servicesQueryRepo repository.ServicesQueryRepo,
) ([]entity.InstallableService, error) {
	return servicesQueryRepo.GetInstallables()
}
