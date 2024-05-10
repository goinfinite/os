package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
)

func DeleteSslPairVhosts(
	sslQueryRepo repository.SslQueryRepo,
	sslCmdRepo repository.SslCmdRepo,
	dto dto.DeleteSslPairVhosts,
) error {
	_, err := sslQueryRepo.ReadById(dto.SslPairId)
	if err != nil {
		return errors.New("SslPairNotFound")
	}

	err = sslCmdRepo.DeleteSslPairVhosts(dto)
	if err != nil {
		log.Printf("DeleteSslPairVhostsError: %s", err.Error())
		return errors.New("DeleteSslPairVhostsInfraError")
	}

	return nil
}
