package useCase

import (
	"errors"
	"log"
	"strings"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
)

func CreateVirtualHost(
	vhostQueryRepo repository.VirtualHostQueryRepo,
	vhostCmdRepo repository.VirtualHostCmdRepo,
	addVirtualHost dto.CreateVirtualHost,
) error {
	_, err := vhostQueryRepo.GetByHostname(addVirtualHost.Hostname)
	if err == nil {
		return errors.New("VirtualHostAlreadyExists")
	}

	isAlias := addVirtualHost.Type.String() == "alias"
	if isAlias && addVirtualHost.ParentHostname == nil {
		return errors.New("AliasMustHaveParentHostname")
	}

	hostnameStr := addVirtualHost.Hostname.String()
	hasWildcardInHostname := strings.HasPrefix(hostnameStr, "*.")
	if hasWildcardInHostname {
		hostnameWithoutWildcardStr := strings.Replace(hostnameStr, "*.", "", 1)
		hostnameWithoutWildcard, err := valueObject.NewFqdn(hostnameWithoutWildcardStr)
		if err != nil {
			return errors.New("FailedToRemoveWildcardFromHostname: " + err.Error())
		}

		addVirtualHost.Hostname = hostnameWithoutWildcard
	}

	err = vhostCmdRepo.Create(addVirtualHost)
	if err != nil {
		log.Printf("CreateVirtualHostError: %s", err.Error())
		return errors.New("CreateVirtualHostInfraError")
	}

	log.Printf("VirtualHost '%s' created.", addVirtualHost.Hostname)

	return nil
}
