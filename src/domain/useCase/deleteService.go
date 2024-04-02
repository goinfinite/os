package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
)

func DeleteService(
	queryRepo repository.ServicesQueryRepo,
	cmdRepo repository.ServicesCmdRepo,
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
