package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
)

func DeleteService(
	servicesQueryRepo repository.ServicesQueryRepo,
	servicesCmdRepo repository.ServicesCmdRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	svcName valueObject.ServiceName,
) error {
	serviceEntity, err := servicesQueryRepo.ReadByName(svcName)
	if err != nil {
		return errors.New("ServiceNotFound")
	}

	isSystemService := serviceEntity.Type.String() == "system"
	if isSystemService {
		return errors.New("SystemServicesCannotBeUninstalled")
	}

	err = mappingCmdRepo.DeleteAuto(svcName)
	if err != nil {
		log.Printf("DeleteAutoMappingError: %s", err.Error())
		return errors.New("DeleteAutoMappingsInfraError")
	}

	err = servicesCmdRepo.Delete(svcName)
	if err != nil {
		log.Printf("DeleteServiceError: %v", err)
		return errors.New("DeleteServiceInfraError")
	}

	log.Printf("Service '%v' deleted.", svcName.String())

	return nil
}
