package useCase

import (
	"errors"
	"log"
	"slices"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
)

func AddSslPair(
	sslCmdRepo repository.SslCmdRepo,
	vhostQueryRepo repository.VirtualHostQueryRepo,
	addSslPair dto.AddSslPair,
) error {
	existingVhosts, err := vhostQueryRepo.Get()
	if err != nil {
		log.Printf("FailedToGetVhosts: %s", err.Error())
		return errors.New("FailedToGetVhostsInfraError")
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

		if slices.Contains(addSslPair.VirtualHosts, vhost.Hostname) {
			validSslVirtualHosts = append(validSslVirtualHosts, vhost.Hostname)
		}
	}

	if len(validSslVirtualHosts) == 0 {
		return errors.New("VhostDoesNotExists")
	}

	addSslPair.VirtualHosts = validSslVirtualHosts

	err = sslCmdRepo.Add(addSslPair)
	if err != nil {
		log.Printf("AddSslPairError: %s", err)
		return errors.New("AddSslPairInfraError")
	}

	return nil
}
