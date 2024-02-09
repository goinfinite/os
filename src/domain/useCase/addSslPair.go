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
	vhosts, err := vhostQueryRepo.Get()
	if err != nil {
		log.Printf("FailedToGetVhosts: %s", err.Error())
		return errors.New("FailedToGetVhostsInfraError")
	}

	if len(vhosts) == 0 {
		log.Printf("VhostsNotFound")
		return errors.New("VhostsNotFound")
	}

	vhostsStr := []string{}
	for _, vhost := range vhosts {
		vhostsStr = append(vhostsStr, vhost.Hostname.String())
	}

	shouldReturnVhostError := false
	if len(addSslPair.VirtualHosts) == 1 {
		shouldReturnVhostError = true
	}

	existingVhosts := []valueObject.Fqdn{}
	for _, vhost := range addSslPair.VirtualHosts {
		if !slices.Contains(vhostsStr, vhost.String()) {
			log.Printf("VhostDoesNotExists: %s", vhost.String())
			if shouldReturnVhostError {
				return errors.New("VhostDoesNotExists")
			}

			continue
		}

		existingVhosts = append(existingVhosts, vhost)
	}

	addSslPair.VirtualHosts = existingVhosts

	err = sslCmdRepo.Add(addSslPair)
	if err != nil {
		log.Printf("AddSslPairError: %s", err)
		return errors.New("AddSslPairInfraError")
	}

	return nil
}
