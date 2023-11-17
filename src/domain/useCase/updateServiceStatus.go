package useCase

import (
	"errors"
	"log"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
)

var servicesDependencies = map[string][]string{
	"php": {"openlitespeed"},
}

func isServiceInstalled(
	servicesQueryRepo repository.ServicesQueryRepo,
	serviceName valueObject.ServiceName,
) bool {
	serviceStatus, err := servicesQueryRepo.GetByName(serviceName)
	if err != nil {
		return false
	}

	return serviceStatus.Status.String() != "uninstalled"
}

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

	servicesWithDependencies := maps.Keys(servicesDependencies)
	serviceHasDependencies := slices.Contains(
		servicesWithDependencies,
		updateSvcStatusDto.Name.String(),
	)
	if serviceHasDependencies {
		serviceDependencies := servicesDependencies[updateSvcStatusDto.Name.String()]

		for _, dependency := range serviceDependencies {
			dependencyServiceName := valueObject.NewServiceNamePanic(dependency)
			isDependencyInstalled := isServiceInstalled(
				servicesQueryRepo,
				dependencyServiceName,
			)

			if isDependencyInstalled {
				continue
			}

			err = servicesCmdRepo.Install(dependencyServiceName, nil)
			if err != nil {
				log.Printf("DependenciesInstallFailed: %s", dependency)
				return errors.New("DependenciesInstallFailed")
			}
		}
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
