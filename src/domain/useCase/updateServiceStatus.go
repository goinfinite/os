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
		return errors.New("ServiceStatusAlreadySet")
	}

	isInstalled := currentSvcStatus.Status.String() == "installed"
	isRunning := currentSvcStatus.Status.String() == "running"
	isStopped := currentSvcStatus.Status.String() == "stopped"
	isInstalled = isInstalled || isRunning || isStopped

	shouldRun := updateSvcStatusDto.Status.String() == "running"
	shouldStop := updateSvcStatusDto.Status.String() == "stopped"
	shouldUninstall := updateSvcStatusDto.Status.String() == "uninstalled"
	if !isInstalled && (shouldRun || shouldStop || shouldUninstall) {
		return errors.New("ServiceNotInstalled")
	}

	shouldInstall := updateSvcStatusDto.Status.String() == "installed"
	if isInstalled && shouldInstall {
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
