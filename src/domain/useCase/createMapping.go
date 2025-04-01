package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
)

func mappingTargetLinter(createDto dto.CreateMapping) (dto.CreateMapping, error) {
	if createDto.TargetType == valueObject.MappingTargetTypeStaticFiles {
		return createDto, nil
	}

	hasTargetValue := createDto.TargetValue != nil
	hasResponseCode := createDto.TargetHttpResponseCode != nil
	if !hasTargetValue && !hasResponseCode {
		return createDto, errors.New("MappingMustHaveValueOrResponseCode")
	}

	isResponseCodeTarget := createDto.TargetType == valueObject.MappingTargetTypeResponseCode
	if !isResponseCodeTarget && !hasTargetValue {
		return createDto, errors.New("MappingMustHaveTargetValue")
	}

	if isResponseCodeTarget && hasTargetValue {
		targetValueStr := createDto.TargetValue.String()
		httpRespondeCode, err := valueObject.NewHttpResponseCode(targetValueStr)
		if err != nil {
			return createDto, errors.New("MappingResponseCodeInvalid")
		}
		createDto.TargetHttpResponseCode = &httpRespondeCode
		createDto.TargetValue = nil
	}

	defaultResponseCode, _ := valueObject.NewHttpResponseCode(200)
	if createDto.TargetType == valueObject.MappingTargetTypeUrl {
		defaultResponseCode, _ = valueObject.NewHttpResponseCode(301)
	}
	if !hasResponseCode && createDto.TargetType != valueObject.MappingTargetTypeService {
		createDto.TargetHttpResponseCode = &defaultResponseCode
	}

	return createDto, nil
}

func CreateMapping(
	vhostQueryRepo repository.VirtualHostQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	svcsQueryRepo repository.ServicesQueryRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	createDto dto.CreateMapping,
) error {
	withMappings := true
	virtualHostWithMappings, err := vhostQueryRepo.ReadFirstWithMappings(dto.ReadVirtualHostsRequest{
		Hostname:     &createDto.Hostname,
		WithMappings: &withMappings,
	})
	if err != nil {
		slog.Error("ReadFirstVirtualHostWithMappingsError", slog.String("err", err.Error()))
		return errors.New("ReadFirstVirtualHostWithMappingsInfraError")
	}

	if virtualHostWithMappings.Type == valueObject.VirtualHostTypeAlias {
		return errors.New("AliasCannotHaveMappings")
	}

	createDto, err = mappingTargetLinter(createDto)
	if err != nil {
		slog.Error("MappingTargetLinterError", slog.String("err", err.Error()))
		return errors.New("MappingTargetLinterError")
	}

	for _, mappingEntity := range virtualHostWithMappings.Mappings {
		if mappingEntity.MatchPattern != createDto.MatchPattern {
			continue
		}

		if mappingEntity.Path != createDto.Path {
			continue
		}

		return errors.New("MappingAlreadyExists")
	}

	if createDto.TargetType == valueObject.MappingTargetTypeService {
		serviceName, err := valueObject.NewServiceName(createDto.TargetValue.String())
		if err != nil {
			return errors.New("MappingTargetServiceNameInvalid")
		}

		serviceEntity, err := svcsQueryRepo.ReadFirstInstalledItem(
			dto.ReadFirstInstalledServiceItemsRequest{ServiceName: &serviceName},
		)
		if err != nil {
			slog.Error("ReadServiceEntityError", slog.String("err", err.Error()))
			return errors.New("ReadServiceEntityInfraError")
		}

		if len(serviceEntity.PortBindings) == 0 {
			return errors.New("ServiceDoesNotExposeAnyPorts")
		}

		for _, portBinding := range serviceEntity.PortBindings {
			protocolStr := portBinding.Protocol.String()
			isTransportLayer := protocolStr == "tcp" || protocolStr == "udp"
			if isTransportLayer {
				return errors.New("TransportLayerMappingNotSupportedYet")
			}
		}
	}

	mappingId, err := mappingCmdRepo.Create(createDto)
	if err != nil {
		slog.Error("CreateMappingError", slog.String("err", err.Error()))
		return errors.New("CreateMappingInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		CreateMapping(createDto, mappingId)

	return nil
}
