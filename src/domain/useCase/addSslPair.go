package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
)

func AddSslPair(
	sslCmdRepo repository.SslCmdRepo,
	vhostQueryRepo repository.VirtualHostQueryRepo,
	addSslPair dto.AddSslPair,
) error {
	for _, vhost := range addSslPair.VirtualHosts {
		_, err := vhostQueryRepo.GetByHostname(vhost)
		if err != nil {
			log.Printf("OneOfTheVhostsDoesNotExists: %s", vhost.String())
			return errors.New("OneOfTheVhostsDoesNotExists")
		}
	}

	err := sslCmdRepo.Add(addSslPair)
	if err != nil {
		log.Printf("AddSslPairError: %s", err)
		return errors.New("AddSslPairInfraError")
	}

	return nil
}
