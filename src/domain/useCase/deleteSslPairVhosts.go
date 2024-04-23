package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
)

func DeleteSslPairVhosts(
	sslQueryRepo repository.SslQueryRepo,
	sslCmdRepo repository.SslCmdRepo,
	vhostQueryRepo repository.VirtualHostQueryRepo,
	deleteSslPairVhosts dto.DeleteSslPairVhosts,
) error {
	_, err := sslQueryRepo.GetSslPairById(deleteSslPairVhosts.SslPairId)
	if err != nil {
		return errors.New("SslPairNotFound")
	}

	for _, vhost := range deleteSslPairVhosts.VirtualHosts {
		_, err := vhostQueryRepo.GetByHostname(vhost)
		if err != nil {
			log.Printf("VhostNotFound: %s", vhost.String())
			continue
		}

		err = sslCmdRepo.RemoveVhostFromSslPair(vhost)
		if err != nil {
			log.Printf(
				"RemoveVhostFromSslPairError (%s): %s",
				vhost.String(),
				err.Error(),
			)
		}
	}

	return nil
}
