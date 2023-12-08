package useCase

import (
	"errors"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/repository"
)

func updateServiceStatus(
	servicesCmdRepo repository.ServicesCmdRepo,
	serviceEntity entity.Service,
	updateDto dto.UpdateService,
) error {
	if serviceEntity.Status.String() == updateDto.Status.String() {
		return errors.New("ServiceStatusAlreadySet")
	}

	isInstalled := serviceEntity.Status.String() != "uninstalled"
	if !isInstalled {
		return errors.New("ServiceNotInstalled")
	}

	isSystemService := serviceEntity.Type.String() == "system"
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
	serviceEntity, err := servicesQueryRepo.GetByName(updateDto.Name)
	if err != nil {
		return err
	}

	if updateDto.Status != nil {
		err = updateServiceStatus(
			servicesCmdRepo,
			serviceEntity,
			updateDto,
		)
		if err != nil {
			return err
		}
	}

	isSoloService := serviceEntity.Type.String() == "solo"
	shouldUpdateVersion := updateDto.Version != nil
	if isSoloService && shouldUpdateVersion {
		return errors.New("SoloServicesVersionCannotBeChanged")
	}

	return servicesCmdRepo.Update(updateDto)
}
