package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
)

func CreateMapping(
	queryRepo repository.VirtualHostQueryRepo,
	cmdRepo repository.VirtualHostCmdRepo,
	svcsQueryRepo repository.ServicesQueryRepo,
	addMapping dto.CreateMapping,
) error {
	vhostWithMappings, err := queryRepo.GetWithMappings()
	if err != nil {
		log.Printf("GetVirtualHostsError: %s", err.Error())
		return errors.New("GetVirtualHostsInfraError")
	}

	vhostIndex := -1
	for vhostWithMappingIndex, vhostWithMapping := range vhostWithMappings {
		if vhostWithMapping.Hostname != addMapping.Hostname {
			continue
		}

		for _, mapping := range vhostWithMapping.Mappings {
			if mapping.MatchPattern != addMapping.MatchPattern {
				continue
			}

			if mapping.Path != addMapping.Path {
				continue
			}

			return errors.New("MappingAlreadyExists")
		}

		vhostIndex = vhostWithMappingIndex
	}

	if vhostIndex == -1 {
		return errors.New("VirtualHostNotFound")
	}

	if vhostWithMappings[vhostIndex].Type.String() == "alias" {
		return errors.New("AliasCannotHaveMappings")
	}

	isServiceTarget := addMapping.TargetType.String() == "service"
	if isServiceTarget {
		if addMapping.TargetServiceName == nil {
			return errors.New("TargetServiceNameRequired")
		}

		service, err := svcsQueryRepo.GetByName(*addMapping.TargetServiceName)
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

	isUrlTarget := addMapping.TargetType.String() == "url"
	if isUrlTarget && addMapping.TargetUrl == nil {
		return errors.New("TargetUrlRequired")
	}

	defaultResponseCode, _ := valueObject.NewHttpResponseCode(301)
	if isUrlTarget && addMapping.TargetHttpResponseCode == nil {
		addMapping.TargetHttpResponseCode = &defaultResponseCode
	}

	isResponseCodeTarget := addMapping.TargetType.String() == "response-code"
	if isResponseCodeTarget && addMapping.TargetHttpResponseCode == nil {
		return errors.New("TargetHttpResponseCodeRequired")
	}

	err = cmdRepo.CreateMapping(addMapping)
	if err != nil {
		log.Printf("CreateMappingError: %s", err.Error())
		return errors.New("CreateMappingInfraError")
	}

	return nil
}
