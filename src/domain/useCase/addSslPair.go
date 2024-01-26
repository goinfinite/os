package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
)

func AddSslPair(
	sslCmdRepo repository.SslCmdRepo,
	addSslPair dto.AddSslPair,
) error {
	err := sslCmdRepo.Add(addSslPair)
	if err != nil {
		log.Printf("AddSslPairError: %s", err)
		return errors.New("AddSslPairInfraError")
	}

	log.Printf(
		"SSL '%v' added in '%v' virtual host.",
		addSslPair.Certificate.Id.String(),
		addSslPair.VirtualHost.String(),
	)

	return nil
}
