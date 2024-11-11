package useCase

import (
	"errors"
	"log"

	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
)

func DeleteSslPair(
	sslQueryRepo repository.SslQueryRepo,
	sslCmdRepo repository.SslCmdRepo,
	sslPairId valueObject.SslPairId,
) error {
	_, err := sslQueryRepo.ReadById(sslPairId)
	if err != nil {
		log.Printf("SslPairNotFound: %s", err)
		return errors.New("SslPairNotFound")
	}

	err = sslCmdRepo.Delete(sslPairId)
	if err != nil {
		log.Printf("DeleteSslPairError: %s", err)
		return errors.New("DeleteSslPairInfraError")
	}

	return nil
}
