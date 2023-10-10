package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/sam/src/domain/repository"
	"github.com/speedianet/sam/src/domain/valueObject"
)

func DeleteSsl(
	sslQueryRepo repository.SslQueryRepo,
	sslCmdRepo repository.SslCmdRepo,
	sslSerialNumber valueObject.SslSerialNumber,
) error {
	_, err := sslQueryRepo.GetById(sslSerialNumber)
	if err != nil {
		log.Printf("SslNotFound: %s", err)
		return errors.New("SslNotFound")
	}

	err = sslCmdRepo.Delete(sslSerialNumber)
	if err != nil {
		log.Printf("DeleteSslError: %s", err)
		return errors.New("DeleteSslInfraError")
	}

	log.Printf("SslSerialNumber '%v' deleted.", sslSerialNumber)

	return nil
}
