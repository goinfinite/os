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
	sslHashId valueObject.SslHashId,
) error {
	sslPair, err := sslQueryRepo.GetSslPairByHashId(sslHashId)
	if err != nil {
		log.Printf("SslPairNotFound: %s", err)
		return errors.New("SslPairNotFound")
	}

	err = sslCmdRepo.Delete(sslHashId)
	if err != nil {
		log.Printf("DeleteSslPairError: %s", err)
		return errors.New("DeleteSslPairInfraError")
	}

	log.Printf(
		"SSL '%v' of '%v' hostname deleted.",
		sslHashId,
		sslPair.Hostname.String(),
	)

	return nil
}
