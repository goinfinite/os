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
		return DeleteService(queryRepo, cmdRepo, updateDto.Name)
	default:
		return errors.New("UnknownServiceStatus")
	}
}

func UpdateService(
	queryRepo repository.ServicesQueryRepo,
	cmdRepo repository.ServicesCmdRepo,
	vhostQueryRepo repository.VirtualHostQueryRepo,
	vhostCmdRepo repository.VirtualHostCmdRepo,
	updateDto dto.UpdateService,
) error {
	serviceEntity, err := queryRepo.GetByName(updateDto.Name)
	if err != nil {
		return err
	}

	if updateDto.Status != nil {
		return updateServiceStatus(
			queryRepo,
			cmdRepo,
			serviceEntity,
			updateDto,
		)
	}

	isSoloService := serviceEntity.Type.String() == "solo"
	shouldUpdateVersion := updateDto.Version != nil
	if isSoloService && shouldUpdateVersion {
		return errors.New("SoloServicesVersionCannotBeChanged")
	}

	err = cmdRepo.Update(updateDto)
	if err != nil {
		log.Printf("UpdateServiceError: %s", err.Error())
		return errors.New("UpdateServiceInfraError")
	}

	if len(updateDto.PortBindings) > 0 {
		vhostsWithMappings, err := vhostQueryRepo.GetWithMappings()
		if err != nil {
			return err
		}

		var mappingsToRecreate []entity.Mapping
		for _, vhostWithMapping := range vhostsWithMappings {
			for _, vhostMapping := range vhostWithMapping.Mappings {
				if vhostMapping.TargetServiceName.String() != updateDto.Name.String() {
					continue
				}

				mappingsToRecreate = append(mappingsToRecreate, vhostMapping)
				break
			}
		}

		for _, mappingToRecreate := range mappingsToRecreate {
			err = vhostCmdRepo.RecreateMapping(mappingToRecreate)
			if err != nil {
				log.Printf("RecreateMappingError: %s", err.Error())
			}
		}
	}

	return nil
}
