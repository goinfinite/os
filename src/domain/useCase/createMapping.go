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

	mappings, err := mappingQueryRepo.GetByHostname(createMapping.Hostname)
	if err != nil {
		return errors.New("MappingsNotFound")
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

	isServiceTarget := createMapping.TargetType.String() == "service"
	if isServiceTarget {
		if createMapping.TargetServiceName == nil {
			return errors.New("TargetServiceNameRequired")
		}

		service, err := svcsQueryRepo.GetByName(*createMapping.TargetServiceName)
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

	isUrlTarget := createMapping.TargetType.String() == "url"
	if isUrlTarget && createMapping.TargetUrl == nil {
		return errors.New("TargetUrlRequired")
	}

	defaultUrlResponseCode, _ := valueObject.NewHttpResponseCode(301)
	if isUrlTarget && createMapping.TargetHttpResponseCode == nil {
		createMapping.TargetHttpResponseCode = &defaultUrlResponseCode
	}

	isTargetHttpResponseCodeMissing := createMapping.TargetHttpResponseCode == nil

	isResponseCodeTarget := createMapping.TargetType.String() == "response-code"
	if isResponseCodeTarget && isTargetHttpResponseCodeMissing {
		return errors.New("TargetHttpResponseCodeRequired")
	}

	isInlineHtmlTarget := createMapping.TargetType.String() == "inline-html"
	if isInlineHtmlTarget {
		if createMapping.TargetInlineHtmlContent == nil {
			return errors.New("TargetInlineHtmlContentRequired")
		}

		if isTargetHttpResponseCodeMissing {
			defaultHttpResponseCode, _ := valueObject.NewHttpResponseCode(200)
			createMapping.TargetHttpResponseCode = &defaultHttpResponseCode
		}
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
