package useCase

import (
	"errors"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
)

func updateServiceStatus(
	servicesQueryRepo repository.ServicesQueryRepo,
	servicesCmdRepo repository.ServicesCmdRepo,
	updateDto dto.UpdateService,
) error {
	currentSvcStatus, err := servicesQueryRepo.GetByName(updateDto.Name)
	if err != nil {
		return err
	}

	if currentSvcStatus.Status.String() == updateDto.Status.String() {
		return errors.New("ServiceStatusAlreadySet")
	}

	isInstalled := currentSvcStatus.Status.String() != "uninstalled"
	if !isInstalled {
		return errors.New("ServiceNotInstalled")
	}

	isSystemService := currentSvcStatus.Type.String() == "system"
	shouldUninstall := updateDto.Status.String() == "uninstalled"
	if isSystemService && shouldUninstall {
		return errors.New("SystemServicesCannotBeUninstalled")
	}

	switch updateDto.Status.String() {
	case "running":
		return servicesCmdRepo.Start(updateDto.Name)
	case "stopped":
		return servicesCmdRepo.Stop(updateDto.Name)
	case "uninstalled":
		return servicesCmdRepo.Uninstall(updateDto.Name)
	default:
		return errors.New("UnknownServiceStatus")
	}
}

func UpdateService(
	servicesQueryRepo repository.ServicesQueryRepo,
	servicesCmdRepo repository.ServicesCmdRepo,
	updateDto dto.UpdateService,
) error {
	var err error
	if updateDto.Status != nil {
		err = updateServiceStatus(
			servicesQueryRepo,
			servicesCmdRepo,
			updateDto,
		)
		if err != nil {
			return err
		}
	}

	return servicesCmdRepo.Update(updateDto)
}
