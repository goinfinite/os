package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/sam/src/domain/repository"
	"github.com/speedianet/sam/src/domain/valueObject"
)

func DeleteSslPair(
	sslQueryRepo repository.SslQueryRepo,
	sslCmdRepo repository.SslCmdRepo,
	sslId valueObject.SslId,
) error {
	sslPair, err := sslQueryRepo.GetSslPairById(sslId)
	if err != nil {
		log.Printf("SslPairNotFound: %s", err)
		return errors.New("SslPairNotFound")
	}

	err = sslCmdRepo.Delete(sslId)
	if err != nil {
		log.Printf("DeleteSslPairError: %s", err)
		return errors.New("DeleteSslPairInfraError")
	}

	log.Printf(
		"SSL '%v' of '%v' hostname deleted.",
		sslId,
		sslPair.Hostname.String(),
	)

	return nil
}
