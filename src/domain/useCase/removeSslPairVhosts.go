package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
)

func RemoveSslPairVhosts(
	sslQueryRepo repository.SslQueryRepo,
	sslCmdRepo repository.SslCmdRepo,
	vhostQueryRepo repository.VirtualHostQueryRepo,
	removeSslPairVhosts dto.RemoveSslPairVhosts,
) error {
	_, err := sslQueryRepo.GetSslPairById(removeSslPairVhosts.SslPairId)
	if err != nil {
		return errors.New("SslPairNotFound")
	}

	for _, vhost := range removeSslPairVhosts.VirtualHosts {
		_, err := vhostQueryRepo.GetByHostname(vhost)
		if err != nil {
			log.Printf("VhostNotFound: %s", vhost.String())
			continue
		}

		err = sslCmdRepo.RemoveVhostFromSslPair(vhost)
		if err != nil {
			log.Printf("RemoveVhostFromSslPairError: %s", err.Error())
		}
	}

	return nil
}
