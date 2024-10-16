package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
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
		slog.Error("DeleteAutoMappingError", slog.Any("error", err))
		return errors.New("DeleteAutoMappingsInfraError")
	}

	err = servicesCmdRepo.Delete(svcName)
	if err != nil {
		slog.Error("DeleteServiceError", slog.Any("error", err))
		return errors.New("DeleteServiceInfraError")
	}

	slog.Info("Service "+svcName.String()+" deleted.", slog.Any("error", err))

	return nil
}
