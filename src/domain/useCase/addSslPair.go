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

	existingVhostsStr := []string{}
	for _, vhost := range existingVhosts {
		existingVhostsStr = append(existingVhostsStr, vhost.Hostname.String())
	}

	validSslVirtualHosts := []valueObject.Fqdn{}
	for _, vhost := range addSslPair.VirtualHosts {
		if !slices.Contains(existingVhostsStr, vhost.String()) {
			log.Printf("VhostDoesNotExists: %s", vhost.String())
			continue
		}

		validSslVirtualHosts = append(validSslVirtualHosts, vhost)
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
