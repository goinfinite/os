package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
)

type DeleteService struct {
	queryRepo      repository.ServicesQueryRepo
	cmdRepo        repository.ServicesCmdRepo
	vhostQueryRepo repository.VirtualHostQueryRepo
	vhostCmdRepo   repository.VirtualHostCmdRepo
}

func NewDeleteService(
	queryRepo repository.ServicesQueryRepo,
	cmdRepo repository.ServicesCmdRepo,
	vhostQueryRepo repository.VirtualHostQueryRepo,
	vhostCmdRepo repository.VirtualHostCmdRepo,
) DeleteService {
	return DeleteService{
		queryRepo:      queryRepo,
		cmdRepo:        cmdRepo,
		vhostQueryRepo: vhostQueryRepo,
		vhostCmdRepo:   vhostCmdRepo,
	}
}

func (uc DeleteService) deleteServiceMapping(
	svcName valueObject.ServiceName,
) error {
	vhostsWithMappings, err := uc.vhostQueryRepo.GetWithMappings()
	if err != nil {
		return err
	}

	if len(vhostsWithMappings) == 0 {
		return nil
	}

	primaryVhostWithMapping := vhostsWithMappings[0]
	if len(primaryVhostWithMapping.Mappings) == 0 {
		return nil
	}

	var mappingToDelete entity.Mapping
	for _, primaryVhostMapping := range primaryVhostWithMapping.Mappings {
		if primaryVhostMapping.TargetType.String() != "service" {
			continue
		}

		targetServiceName := primaryVhostMapping.TargetServiceName
		if targetServiceName == nil {
			continue
		}

		if targetServiceName.String() != svcName.String() {
			continue
		}

		mappingToDelete = primaryVhostMapping
	}

	hasMappingToDelete := mappingToDelete.Hostname != ""
	if !hasMappingToDelete {
		return nil
	}

	err = uc.vhostCmdRepo.DeleteMapping(mappingToDelete)
	if err != nil {
		return err
	}

	return nil
}

func (uc DeleteService) Execute(svcName valueObject.ServiceName) error {
	serviceEntity, err := uc.queryRepo.GetByName(svcName)
	if err != nil {
		return errors.New("ServiceNotFound")
	}

	isSystemService := serviceEntity.Type.String() == "system"
	if isSystemService {
		return errors.New("SystemServicesCannotBeUninstalled")
	}

	err = uc.deleteServiceMapping(svcName)
	if err != nil {
		log.Printf("DeleteServiceMappingError: %s", err.Error())
		return errors.New("DeleteServiceMappingsInfraError")
	}

	err = uc.cmdRepo.Uninstall(svcName)
	if err != nil {
		log.Printf("DeleteServiceError: %v", err)
		return errors.New("DeleteServiceInfraError")
	}

	log.Printf("Service '%v' deleted.", svcName.String())

	return nil
}
