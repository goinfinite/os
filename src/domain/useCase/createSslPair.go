package useCase

import (
	"errors"
	"log"
	"slices"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
)

func CreateSslPair(
	sslCmdRepo repository.SslCmdRepo,
	vhostQueryRepo repository.VirtualHostQueryRepo,
	createSslPair dto.CreateSslPair,
) error {
	existingVhosts, err := vhostQueryRepo.Read()
	if err != nil {
		log.Printf("ReadVhostsError: %s", err.Error())
		return errors.New("ReadVhostsInfraError")
	}

	if len(existingVhosts) == 0 {
		log.Printf("VhostsNotFound")
		return errors.New("VhostsNotFound")
	}

	validSslVirtualHostsHostnames := []valueObject.Fqdn{}
	for _, vhost := range existingVhosts {
		if vhost.Type.String() == "alias" {
			continue
		}

		if slices.Contains(createSslPair.VirtualHostsHostnames, vhost.Hostname) {
			validSslVirtualHostsHostnames = append(
				validSslVirtualHostsHostnames, vhost.Hostname,
			)
		}
	}

	if len(validSslVirtualHostsHostnames) == 0 {
		return errors.New("VhostDoesNotExists")
	}

	createSslPair.VirtualHostsHostnames = validSslVirtualHostsHostnames

	err = sslCmdRepo.Create(createSslPair)
	if err != nil {
		log.Printf("CreateSslPairError: %s", err)
		return errors.New("CreateSslPairInfraError")
	}

	return nil
}
