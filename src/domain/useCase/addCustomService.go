package useCase

import (
	"errors"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
)

func AddCustomService(
	servicesQueryRepo repository.ServicesQueryRepo,
	servicesCmdRepo repository.ServicesCmdRepo,
	addDto dto.AddCustomService,
) error {
	currentSvcStatus, err := servicesQueryRepo.GetByName(addDto.Name)
	if err != nil {
		return err
	}

	isInstalled := currentSvcStatus.Status.String() != "uninstalled"
	if isInstalled {
		return errors.New("ServiceAlreadyInstalled")
	}

	return servicesCmdRepo.AddCustom(addDto)
}
