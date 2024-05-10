package useCase

import (
	"errors"
	"log"

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

	isResponseCodeMapping := targetTypeStr == "response-code"
	isTargetValueRequired := !isResponseCodeMapping
	if isTargetValueRequired && !hasTargetValue {
		return errors.New("MappingMustHaveValue")
	}

	if isResponseCodeMapping && !hasTargetHttpResponseCode {
		targetValuetr := createMapping.TargetValue.String()
		httpRespondeCode, err := valueObject.NewHttpResponseCode(targetValuetr)
		if err != nil {
			return err
		}

		createMapping.TargetHttpResponseCode = &httpRespondeCode
		createMapping.TargetValue = nil
	}

	mappings, err := mappingQueryRepo.ReadByHostname(createMapping.Hostname)
	if err != nil {
		log.Printf("ReadMappingsError: %s", err.Error())
		return errors.New("ReadMappingsInfraError")
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
			log.Printf("GetServiceByNameError: %s", err.Error())
			return errors.New("GetServiceByNameInfraError")
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

	_, err = mappingCmdRepo.Create(createMapping)
	if err != nil {
		log.Printf("CreateMappingError: %s", err.Error())
		return errors.New("CreateMappingInfraError")
	}

	return nil
}
