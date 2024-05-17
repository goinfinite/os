package useCase

import (
	"errors"
	"log"
	"slices"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
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

	validSslVirtualHosts := []valueObject.Fqdn{}
	for _, vhost := range existingVhosts {
		if vhost.Type.String() == "alias" {
			continue
		}

		if slices.Contains(createSslPair.VirtualHosts, vhost.Hostname) {
			validSslVirtualHosts = append(validSslVirtualHosts, vhost.Hostname)
		}
	}

	if len(validSslVirtualHosts) == 0 {
		return errors.New("VhostDoesNotExists")
	}

	createSslPair.VirtualHosts = validSslVirtualHosts

	err = sslCmdRepo.Create(createSslPair)
	if err != nil {
		log.Printf("CreateSslPairError: %s", err)
		return errors.New("CreateSslPairInfraError")
	}

	return nil
}
