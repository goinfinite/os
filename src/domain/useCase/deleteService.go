package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
)

/*func (uc DeleteService) deleteSvcAutoMapping(
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
}*/

func DeleteService(
	queryRepo repository.ServicesQueryRepo,
	cmdRepo repository.ServicesCmdRepo,
	vhostQueryRepo repository.VirtualHostQueryRepo,
	vhostCmdRepo repository.VirtualHostCmdRepo,
	svcName valueObject.ServiceName,
) error {
	serviceEntity, err := queryRepo.GetByName(svcName)
	if err != nil {
		return errors.New("ServiceNotFound")
	}

	isSystemService := serviceEntity.Type.String() == "system"
	if isSystemService {
		return errors.New("SystemServicesCannotBeUninstalled")
	}

	err = vhostCmdRepo.DeleteAutoMapping(svcName)
	if err != nil {
		log.Printf("DeleteAutoMappingError: %s", err.Error())
		return errors.New("DeleteAutoMappingsInfraError")
	}

	err = cmdRepo.Uninstall(svcName)
	if err != nil {
		log.Printf("DeleteServiceError: %v", err)
		return errors.New("DeleteServiceInfraError")
	}

	log.Printf("Service '%v' deleted.", svcName.String())

	return nil
}
