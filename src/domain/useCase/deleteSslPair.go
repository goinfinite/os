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
	sslSerialNumber valueObject.SslSerialNumber,
) error {
	sslPair, err := sslQueryRepo.GetSslPairBySerialNumber(sslSerialNumber)
	if err != nil {
		log.Printf("SslPairNotFound: %s", err)
		return errors.New("SslPairNotFound")
	}

	err = sslCmdRepo.Delete(sslSerialNumber)
	if err != nil {
		log.Printf("DeleteSslPairError: %s", err)
		return errors.New("DeleteSslPairInfraError")
	}

	log.Printf(
		"SSL '%v' deleted in '%v' hostname.",
		sslSerialNumber,
		sslPair.Hostname.String(),
	)

	return nil
}
