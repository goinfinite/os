package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
)

func AddMapping(
	queryRepo repository.VirtualHostQueryRepo,
	cmdRepo repository.VirtualHostCmdRepo,
	addMapping dto.AddMapping,
) error {
	vhostWithMappings, err := queryRepo.GetWithMappings()
	if err != nil {
		log.Printf("GetVirtualHostsError: %s", err.Error())
		return errors.New("GetVirtualHostsInfraError")
	}

	var vhost *entity.VirtualHost
	for _, vhostWithMapping := range vhostWithMappings {
		if vhostWithMapping.Hostname != addMapping.Hostname {
			continue
		}

		vhost = &vhostWithMapping.VirtualHost

		for _, mapping := range vhostWithMapping.Mappings {
			if mapping.MatchPattern != addMapping.MatchPattern {
				continue
			}

			if mapping.Path != addMapping.Path {
				continue
			}

			return errors.New("MappingAlreadyExists")
		}
	}

	if vhost == nil {
		return errors.New("VirtualHostNotFound")
	}

	if vhost.Type.String() == "alias" {
		return errors.New("AliasCannotHaveMappings")
	}

	isServiceTarget := addMapping.TargetType.String() == "service"
	if isServiceTarget && addMapping.TargetService == nil {
		return errors.New("TargetServiceNameRequired")
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

	err = cmdRepo.AddMapping(addMapping)
	if err != nil {
		log.Printf("AddMappingError: %s", err.Error())
		return errors.New("AddMappingInfraError")
	}

	log.Printf("Mapping '%s' added.", addMapping.Hostname)

	return nil
}
