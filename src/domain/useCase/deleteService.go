package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
	tkDto "github.com/goinfinite/tk/src/domain/dto"
	tkRepository "github.com/goinfinite/tk/src/domain/repository"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

func DeleteServiceMappings(
	mappingQueryRepo repository.MappingQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	serviceName valueObject.ServiceName,
	shouldRecreate bool,
	operatorAccountId tkValueObject.AccountId,
	operatorIpAddress tkValueObject.IpAddress,
) error {
	targetType := valueObject.MappingTargetTypeService
	serviceNameStr := serviceName.String()

	targetValue, err := valueObject.NewMappingTargetValue(serviceNameStr, targetType)
	if err != nil {
		return err
	}

	readMappingsResponse, err := mappingQueryRepo.Read(dto.ReadMappingsRequest{
		Pagination:  tkDto.PaginationUnpaginated,
		TargetType:  &targetType,
		TargetValue: &targetValue,
	})
	if err != nil {
		return errors.New("ReadMappingsInfraError: " + err.Error())
	}

	if len(readMappingsResponse.Mappings) == 0 {
		return nil
	}

	for _, mappingEntity := range readMappingsResponse.Mappings {
		err = mappingCmdRepo.Delete(mappingEntity.Id)
		if err != nil {
			slog.Error(
				"DeleteMappingInfraError",
				slog.String("err", err.Error()),
				slog.String("mappingId", mappingEntity.Id.String()),
				slog.String("serviceName", serviceName.String()),
				slog.String("method", "recreateServiceAutoMapping"),
			)
			continue
		}

		if !shouldRecreate {
			continue
		}

		_, err = mappingCmdRepo.Create(dto.NewCreateMapping(
			mappingEntity.Hostname, mappingEntity.Path,
			mappingEntity.MatchPattern, mappingEntity.TargetType,
			mappingEntity.TargetValue, mappingEntity.TargetHttpResponseCode,
			mappingEntity.ShouldUpgradeInsecureRequests, mappingEntity.MappingSecurityRuleId,
			operatorAccountId, operatorIpAddress,
		))
		if err != nil {
			slog.Error(
				"RecreateMappingInfraError",
				slog.String("err", err.Error()),
				slog.String("mappingId", mappingEntity.Id.String()),
				slog.String("serviceName", serviceName.String()),
				slog.String("method", "recreateServiceAutoMapping"),
			)
			continue
		}
	}

	return nil
}

func DeleteService(
	servicesQueryRepo repository.ServicesQueryRepo,
	servicesCmdRepo repository.ServicesCmdRepo,
	mappingQueryRepo repository.MappingQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	activityRecordCmdRepo tkRepository.ActivityRecordCmdRepo,
	deleteDto dto.DeleteService,
) error {
	serviceEntity, err := servicesQueryRepo.ReadFirstInstalledItem(
		dto.ReadFirstInstalledServiceItemsRequest{ServiceName: &deleteDto.Name},
	)
	if err != nil {
		slog.Error("ReadServiceInfraEntityError", slog.String("err", err.Error()))
		return errors.New("ReadServiceEntityError")
	}

	if serviceEntity.Type == valueObject.ServiceTypeSystem {
		return errors.New("SystemServicesCannotBeUninstalled")
	}

	err = DeleteServiceMappings(
		mappingQueryRepo, mappingCmdRepo, deleteDto.Name, false,
		deleteDto.OperatorAccountId, deleteDto.OperatorIpAddress,
	)
	if err != nil {
		slog.Error("DeleteServiceMappingsError", slog.String("err", err.Error()))
		return errors.New("DeleteServiceMappingsInfraError")
	}

	err = servicesCmdRepo.Delete(deleteDto.Name)
	if err != nil {
		slog.Error("DeleteServiceError", slog.String("err", err.Error()))
		return errors.New("DeleteServiceInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).DeleteService(deleteDto)

	return nil
}
