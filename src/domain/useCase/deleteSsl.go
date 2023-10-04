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
	sslId valueObject.SslId,
) error {
	_, err := sslQueryRepo.GetById(sslId)
	if err != nil {
		log.Printf("SslNotFound: %s", err)
		return errors.New("SslNotFound")
	}

	err = sslCmdRepo.Delete(sslId)
	if err != nil {
		log.Printf("DeleteSslError: %s", err)
		return errors.New("DeleteSslInfraError")
	}

	log.Printf("SslId '%v' deleted.", sslId)

	return nil
}
