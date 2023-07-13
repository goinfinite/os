package useCase

import (
	"errors"

	"github.com/speedianet/sam/src/domain/dto"
	"github.com/speedianet/sam/src/domain/repository"
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
		return nil
	}

	isRunning := currentSvcStatus.Status.String() == "running"
	isStopped := currentSvcStatus.Status.String() == "stopped"

	if isRunning || isStopped &&
		updateSvcStatusDto.Status.String() == "installed" {
		return errors.New("ServiceAlreadyInstalled")
	}

	switch updateSvcStatusDto.Status.String() {
	case "running":
		return servicesCmdRepo.Start(updateSvcStatusDto.Name)
	case "stopped":
		return servicesCmdRepo.Stop(updateSvcStatusDto.Name)
	case "installed":
		return servicesCmdRepo.Install(
			updateSvcStatusDto.Name,
			updateSvcStatusDto.Version,
		)
	case "uninstalled":
		return servicesCmdRepo.Uninstall(updateSvcStatusDto.Name)
	default:
		return errors.New("UnknownServiceStatus")
	}
}
