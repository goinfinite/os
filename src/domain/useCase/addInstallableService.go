package useCase

import (
	"errors"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
)

func AddInstallableService(
	servicesQueryRepo repository.ServicesQueryRepo,
	servicesCmdRepo repository.ServicesCmdRepo,
	addDto dto.AddInstallableService,
) error {
	currentSvcStatus, err := servicesQueryRepo.GetByName(addDto.Name)
	if err != nil {
		return err
	}

	isInstalled := currentSvcStatus.Status.String() != "uninstalled"
	if isInstalled {
		return errors.New("ServiceAlreadyInstalled")
	}

	isSystemService := currentSvcStatus.Type.String() == "system"
	if isSystemService {
		return errors.New("SystemServicesCannotBeInstalled")
	}

	return servicesCmdRepo.AddInstallable(addDto)
}
