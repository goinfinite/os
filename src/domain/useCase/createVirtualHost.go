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
	createVirtualHost dto.CreateVirtualHost,
) error {
	_, err := vhostQueryRepo.ReadByHostname(createVirtualHost.Hostname)
	if err == nil {
		return errors.New("VirtualHostAlreadyExists")
	}

	isAlias := createVirtualHost.Type.String() == "alias"
	if isAlias && createVirtualHost.ParentHostname == nil {
		return errors.New("AliasMustHaveParentHostname")
	}

	hostnameStr := createVirtualHost.Hostname.String()
	hasWildcardInHostname := strings.HasPrefix(hostnameStr, "*.")
	if hasWildcardInHostname {
		hostnameWithoutWildcardStr := strings.Replace(hostnameStr, "*.", "", 1)
		hostnameWithoutWildcard, err := valueObject.NewFqdn(hostnameWithoutWildcardStr)
		if err != nil {
			return errors.New("FailedToRemoveWildcardFromHostname: " + err.Error())
		}

		createVirtualHost.Hostname = hostnameWithoutWildcard
	}

	err = vhostCmdRepo.Create(createVirtualHost)
	if err != nil {
		log.Printf("CreateVirtualHostError: %s", err.Error())
		return errors.New("CreateVirtualHostInfraError")
	}

	log.Printf("VirtualHost '%s' created.", createVirtualHost.Hostname)

	return nil
}
