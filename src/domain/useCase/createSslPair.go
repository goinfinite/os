package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
)

func CreateSslPair(
	sslCmdRepo repository.SslCmdRepo,
	addSslPair dto.CreateSslPair,
) error {
	err := sslCmdRepo.Create(addSslPair)
	if err != nil {
		log.Printf("CreateSslPairError: %s", err)
		return errors.New("CreateSslPairInfraError")
	}

	log.Printf(
		"SSL '%v' created in '%v' hostname.",
		addSslPair.Certificate.Id.String(),
		addSslPair.Hostname.String(),
	)

	return nil
}
