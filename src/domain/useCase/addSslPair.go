package useCase

import (
	"errors"
	"log"
	"slices"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
)

func AddSslPair(
	sslCmdRepo repository.SslCmdRepo,
	vhostQueryRepo repository.VirtualHostQueryRepo,
	addSslPair dto.AddSslPair,
) error {
	allVhosts, err := vhostQueryRepo.Get()
	if err != nil {
		log.Printf("FailedToGetVhosts: %s", err.Error())
		return errors.New("FailedToGetVhostsInfraError")
	}

	if len(allVhosts) == 0 {
		log.Printf("VhostsNotFound")
		return errors.New("VhostsNotFound")
	}

	existingVhosts := []string{}
	for _, vhost := range allVhosts {
		existingVhosts = append(existingVhosts, vhost.Hostname.String())
	}

	shouldReturnVhostError := false
	if len(addSslPair.VirtualHosts) == 1 {
		shouldReturnVhostError = true
	}

	for _, vhost := range addSslPair.VirtualHosts {
		if !slices.Contains(existingVhosts, vhost.String()) {
			log.Printf("VhostDoesNotExists: %s", vhost.String())
			if shouldReturnVhostError {
				return errors.New("VhostDoesNotExists")
			}
		}
	}

	err = sslCmdRepo.Add(addSslPair)
	if err != nil {
		log.Printf("AddSslPairError: %s", err)
		return errors.New("AddSslPairInfraError")
	}

	return nil
}
