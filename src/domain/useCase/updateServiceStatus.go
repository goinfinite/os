package useCase

import (
	"errors"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
)

func UpdateServiceStatus(
	servicesQueryRepo repository.ServicesQueryRepo,
	servicesCmdRepo repository.ServicesCmdRepo,
	updateSvcStatusDto dto.UpdateSvcStatus,
) error {
	currentSvcStatus, err := servicesQueryRepo.GetByName(updateSvcStatusDto.Name)
	if err != nil {
		return err
	}

	if currentSvcStatus.Status.String() == updateSvcStatusDto.Status.String() {
		return errors.New("ServiceStatusAlreadySet")
	}

	isInstalled := currentSvcStatus.Status.String() != "uninstalled"
	if !isInstalled {
		return errors.New("ServiceNotInstalled")
	}

	isSystemService := currentSvcStatus.Type.String() == "system"
	shouldUninstall := updateSvcStatusDto.Status.String() == "uninstalled"
	if isSystemService && shouldUninstall {
		return errors.New("SystemServicesCannotBeUninstalled")
	}

	switch updateSvcStatusDto.Status.String() {
	case "running":
		return servicesCmdRepo.Start(updateSvcStatusDto.Name)
	case "stopped":
		return servicesCmdRepo.Stop(updateSvcStatusDto.Name)
	case "uninstalled":
		return servicesCmdRepo.Uninstall(updateSvcStatusDto.Name)
	default:
		return errors.New("UnknownServiceStatus")
	}
}
