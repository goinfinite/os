package useCase

import (
	"errors"
	"slices"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
	"golang.org/x/exp/maps"
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
	shouldInstall := updateSvcStatusDto.Status.String() == "installed"
	if isInstalled && shouldInstall {
		return errors.New("ServiceAlreadyInstalled")
	}

	if !isInstalled && !shouldInstall {
		return errors.New("ServiceNotInstalled")
	}

	isSystemService := currentSvcStatus.Type.String() == "system"
	shouldUninstall := updateSvcStatusDto.Status.String() == "uninstalled"
	if isSystemService && shouldUninstall {
		return errors.New("SystemServicesCannotBeUninstalled")
	}

	nativeServicesNames := maps.Keys(valueObject.NativeSvcNamesWithAliases)
	isNativeService := slices.Contains(nativeServicesNames, updateSvcStatusDto.Name.String())
	if shouldInstall && !isNativeService {
		return errors.New("NotNativeServiceCannotBeInstalled")
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
