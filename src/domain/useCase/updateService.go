package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/repository"
)

func updateServiceStatus(
	servicesQueryRepo repository.ServicesQueryRepo,
	servicesCmdRepo repository.ServicesCmdRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	serviceEntity entity.Service,
	updateDto dto.UpdateService,
) error {
	if serviceEntity.Status.String() == updateDto.Status.String() {
		return nil
	}

	isInstalled := serviceEntity.Status.String() != "uninstalled"
	if !isInstalled {
		return errors.New("ServiceNotInstalled")
	}

	switch updateDto.Status.String() {
	case "running":
		return servicesCmdRepo.Start(updateDto.Name)
	case "stopped":
		return servicesCmdRepo.Stop(updateDto.Name)
	case "uninstalled":
		return DeleteService(
			servicesQueryRepo, servicesCmdRepo, mappingCmdRepo, updateDto.Name,
		)
	case "restarting":
		return servicesCmdRepo.Restart(updateDto.Name)
	default:
		return errors.New("UnknownServiceStatus")
	}
}

func UpdateService(
	servicesQueryRepo repository.ServicesQueryRepo,
	servicesCmdRepo repository.ServicesCmdRepo,
	mappingQueryRepo repository.MappingQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	updateDto dto.UpdateService,
) error {
	serviceEntity, err := servicesQueryRepo.GetByName(updateDto.Name)
	if err != nil {
		return err
	}

	isSoloService := serviceEntity.Type.String() == "solo"
	shouldUpdateStatus := updateDto.Status != nil
	if isSoloService && !shouldUpdateStatus {
		return errors.New("SoloServicesCanOnlyChangeStatus")
	}

	if shouldUpdateStatus {
		err = updateServiceStatus(
			servicesQueryRepo,
			servicesCmdRepo,
			mappingCmdRepo,
			serviceEntity,
			updateDto,
		)
		if err != nil {
			log.Printf("UpdateServiceStatusError: %s", err.Error())
			return errors.New("UpdateServiceStatusInfraError")
		}
	}

	shouldUpdateType := updateDto.Type != nil
	shouldUpdateCommand := updateDto.Command != nil
	shouldUpdateVersion := updateDto.Version != nil
	shouldUpdateStartupFile := updateDto.StartupFile != nil
	portBindingsChanged := len(updateDto.PortBindings) != 0
	nothingElseChanged := !shouldUpdateType && !shouldUpdateCommand &&
		!shouldUpdateVersion && !shouldUpdateStartupFile && !portBindingsChanged
	if nothingElseChanged {
		return nil
	}

	err = servicesCmdRepo.Update(updateDto)
	if err != nil {
		log.Printf("UpdateServiceError: %s", err.Error())
		return errors.New("UpdateServiceInfraError")
	}

	if len(updateDto.PortBindings) == 0 {
		return nil
	}

	err = mappingCmdRepo.RecreateByServiceName(updateDto.Name)
	if err != nil {
		log.Printf("RecreateMappingError: %s", err.Error())
	}

	return nil
}
