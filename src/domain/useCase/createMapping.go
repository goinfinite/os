package useCase

import (
	"errors"
	"log"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
)

func mappingTargetLinter(createDto dto.CreateMapping) (dto.CreateMapping, error) {
	targetTypeStr := createDto.TargetType.String()
	isStaticFilesMapping := targetTypeStr == "static-files"
	if isStaticFilesMapping {
		return createDto, nil
	}

	hasTargetValue := createDto.TargetValue != nil
	hasTargetHttpResponseCode := createDto.TargetHttpResponseCode != nil
	if !hasTargetValue && !hasTargetHttpResponseCode {
		return createDto, errors.New("MappingMustHaveValueOrResponseCode")
	}

	isResponseCodeMapping := targetTypeStr == "response-code"
	isTargetValueRequired := !isResponseCodeMapping
	if isTargetValueRequired && !hasTargetValue {
		return createDto, errors.New("MappingMustHaveTargetValue")
	}

	if isResponseCodeMapping {
		if hasTargetValue {
			targetValueStr := createDto.TargetValue.String()
			httpRespondeCode, err := valueObject.NewHttpResponseCode(targetValueStr)
			if err != nil {
				return createDto, err
			}

			createDto.TargetHttpResponseCode = &httpRespondeCode
		}
		createDto.TargetValue = nil
	}

	isUrlMapping := targetTypeStr == "url"
	if isUrlMapping && !hasTargetHttpResponseCode {
		targetHttpResponseCode, _ := valueObject.NewHttpResponseCode(301)
		createDto.TargetHttpResponseCode = &targetHttpResponseCode
	}

	isInlineHtmlMapping := targetTypeStr == "inline-html"
	if isInlineHtmlMapping && !hasTargetHttpResponseCode {
		targetHttpResponseCode, _ := valueObject.NewHttpResponseCode(200)
		createDto.TargetHttpResponseCode = &targetHttpResponseCode
	}

	return createDto, nil
}

func CreateMapping(
	mappingQueryRepo repository.MappingQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	vhostQueryRepo repository.VirtualHostQueryRepo,
	svcsQueryRepo repository.ServicesQueryRepo,
	createDto dto.CreateMapping,
) error {
	vhost, err := vhostQueryRepo.ReadByHostname(createDto.Hostname)
	if err != nil {
		return errors.New("VirtualHostNotFound")
	}

	if vhost.Type.String() == "alias" {
		return errors.New("AliasCannotHaveMappings")
	}

	createDto, err = mappingTargetLinter(createDto)
	if err != nil {
		return err
	}

	existingMappings, err := mappingQueryRepo.ReadByHostname(createDto.Hostname)
	if err != nil {
		log.Printf("ReadMappingsError: %s", err.Error())
		return errors.New("ReadMappingsInfraError")
	}

	for _, mapping := range existingMappings {
		if mapping.MatchPattern != createDto.MatchPattern {
			continue
		}

		if mapping.Path != createDto.Path {
			continue
		}

		return errors.New("MappingAlreadyExists")
	}

	targetTypeStr := createDto.TargetType.String()
	if targetTypeStr == "service" {
		targetValueStr := createDto.TargetValue.String()
		svcName, err := valueObject.NewServiceName(targetValueStr)
		if err != nil {
			return errors.New(err.Error() + ": " + targetValueStr)
		}

		readInstalledDto := dto.ReadInstalledServicesItemsRequest{
			Name:                 &svcName,
			ShouldIncludeMetrics: false,
		}
		serviceEntity, err := svcsQueryRepo.ReadUniqueInstalledItem(readInstalledDto)
		if err != nil {
			return err
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

	_, err = mappingCmdRepo.Create(createDto)
	if err != nil {
		log.Printf("CreateMappingError: %s", err.Error())
		return errors.New("CreateMappingInfraError")
	}

	return nil
}
