package useCase

import (
	"errors"
	"log"
	"strings"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
)

func CreateMapping(
	mappingQueryRepo repository.MappingQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	vhostQueryRepo repository.VirtualHostQueryRepo,
	svcsQueryRepo repository.ServicesQueryRepo,
	createMapping dto.CreateMapping,
) error {
	vhost, err := vhostQueryRepo.GetByHostname(createMapping.Hostname)
	if err != nil {
		return errors.New("VhostNotFound")
	}

	if vhost.Type.String() == "alias" {
		return errors.New("AliasCannotHaveMappings")
	}

	hasTargetValue := createMapping.TargetValue != nil
	hasTargetHttpResponseCode := createMapping.TargetHttpResponseCode != nil

	if !hasTargetValue && !hasTargetHttpResponseCode {
		return errors.New("MappingMustHaveValueOrHttpResponseCode")
	}

	targetTypeStr := createMapping.TargetType.String()
	if targetTypeStr == "response-code" {
		if !hasTargetHttpResponseCode {
			httpRespondeCode, _ := valueObject.NewHttpResponseCode(
				createMapping.TargetValue.String(),
			)
			createMapping.TargetHttpResponseCode = &httpRespondeCode
			createMapping.TargetValue = nil
		}
	}

	if !hasTargetValue {
		return errors.New("MappingMustHaveValue")
	}

	mappings, err := mappingQueryRepo.GetByHostname(createMapping.Hostname)
	if err != nil {
		log.Printf("GetMappingsError: %s", err.Error())
		return errors.New("GetMappingsInfraError")
	}

	for _, mapping := range mappings {
		if mapping.MatchPattern != createMapping.MatchPattern {
			continue
		}

		if mapping.Path != createMapping.Path {
			continue
		}

		return errors.New("MappingAlreadyExists")
	}

	if targetTypeStr == "service" {
		targetValueStr := createMapping.TargetValue.String()
		svcName, err := valueObject.NewServiceName(targetValueStr)
		if err != nil {
			return errors.New(err.Error() + ": " + targetValueStr)
		}

		service, err := svcsQueryRepo.GetByName(svcName)
		if err != nil {
			return err
		}

		if len(service.PortBindings) == 0 {
			return errors.New("ServiceDoesNotExposeAnyPorts")
		}

		for _, portBinding := range service.PortBindings {
			isTcp := portBinding.Protocol.String() == "tcp"
			isUdp := portBinding.Protocol.String() == "udp"
			if isTcp || isUdp {
				return errors.New("Layer4MappingNotSupportedYet")
			}
		}
	}

	if targetTypeStr == "url" && !hasTargetHttpResponseCode {
		targetHttpResponseCode, _ := valueObject.NewHttpResponseCode(301)
		createMapping.TargetHttpResponseCode = &targetHttpResponseCode
	}

	if targetTypeStr == "inline-html" && !hasTargetHttpResponseCode {
		targetHttpResponseCode, _ := valueObject.NewHttpResponseCode(200)
		createMapping.TargetHttpResponseCode = &targetHttpResponseCode
	}

	pathStr := createMapping.Path.String()
	pathStartsWithSlash := strings.HasPrefix(pathStr, "/")
	if !pathStartsWithSlash {
		createMapping.Path, err = valueObject.NewMappingPath("/" + pathStr)
		if err != nil {
			return errors.New("CorrectAutoMappingPathError")
		}
	}

	_, err = mappingCmdRepo.Create(createMapping)
	if err != nil {
		log.Printf("CreateMappingError: %s", err.Error())
		return errors.New("CreateMappingInfraError")
	}

	return nil
}
