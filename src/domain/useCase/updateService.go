package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/repository"
)

func updateServiceStatus(
	queryRepo repository.ServicesQueryRepo,
	cmdRepo repository.ServicesCmdRepo,
	mappingCmdRepo repository.MappingCmdRepo,
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

	switch updateDto.Status.String() {
	case "running":
		return cmdRepo.Start(updateDto.Name)
	case "stopped":
		return cmdRepo.Stop(updateDto.Name)
	case "uninstalled":
		return DeleteService(
			queryRepo,
			cmdRepo,
			mappingCmdRepo,
			updateDto.Name,
		)
	default:
		return errors.New("UnknownServiceStatus")
	}
}

func UpdateService(
	queryRepo repository.ServicesQueryRepo,
	cmdRepo repository.ServicesCmdRepo,
	mappingQueryRepo repository.MappingQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	updateDto dto.UpdateService,
) error {
	serviceEntity, err := queryRepo.GetByName(updateDto.Name)
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
			queryRepo,
			cmdRepo,
			mappingCmdRepo,
			serviceEntity,
			updateDto,
		)
		if err != nil {
			log.Printf("UpdateServiceStatusError: %s", err.Error())
			return errors.New("UpdateServiceStatusInfraError")
		}
	}

	err = cmdRepo.Update(updateDto)
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
