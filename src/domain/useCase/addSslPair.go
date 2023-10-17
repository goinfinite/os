package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/sam/src/domain/dto"
	"github.com/speedianet/sam/src/domain/repository"
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
		"SSL '%v' added in '%v' hostname.",
		addSslPair.Certificate.HashId.String(),
		addSslPair.Hostname.String(),
	)

	return nil
}
