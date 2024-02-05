package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
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

	err = vhostCmdRepo.Add(addVirtualHost)
	if err != nil {
		log.Printf("CreateVirtualHostError: %s", err.Error())
		return errors.New("CreateVirtualHostInfraError")
	}

	log.Printf("VirtualHost '%s' created.", addVirtualHost.Hostname)

	return nil
}
