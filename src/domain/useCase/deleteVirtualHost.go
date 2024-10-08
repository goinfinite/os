package useCase

import (
	"errors"
	"log"

	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
)

func DeleteVirtualHost(
	queryRepo repository.VirtualHostQueryRepo,
	cmdRepo repository.VirtualHostCmdRepo,
	primaryHostname valueObject.Fqdn,
	hostname valueObject.Fqdn,
) error {
	isPrimaryHostname := hostname.String() == primaryHostname.String()
	if isPrimaryHostname {
		return errors.New("PrimaryVirtualHostCannotBeDeleted")
	}

	vhost, err := queryRepo.ReadByHostname(hostname)
	if err != nil {
		return errors.New("VirtualHostNotFound")
	}

	err = cmdRepo.Delete(vhost)
	if err != nil {
		log.Printf("DeleteVirtualHostError: %v", err)
		return errors.New("DeleteVirtualHostInfraError")
	}

	log.Printf("VirtualHost '%v' deleted.", hostname)

	return nil
}
