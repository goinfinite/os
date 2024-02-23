package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
)

func DeleteSslPair(
	sslQueryRepo repository.SslQueryRepo,
	sslCmdRepo repository.SslCmdRepo,
	sslId valueObject.SslId,
) error {
	_, err := sslQueryRepo.GetSslPairById(sslId)
	if err != nil {
		log.Printf("SslPairNotFound: %s", err)
		return errors.New("SslPairNotFound")
	}

	err = sslCmdRepo.Delete(sslId)
	if err != nil {
		log.Printf("DeleteSslPairError: %s", err)
		return errors.New("DeleteSslPairInfraError")
	}

	return nil
}
