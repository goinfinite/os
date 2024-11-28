package useCase

import (
	"errors"
	"log"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

func DeleteSslPair(
	sslQueryRepo repository.SslQueryRepo,
	sslCmdRepo repository.SslCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	deleteDto dto.DeleteSslPair,
) error {
	_, err := sslQueryRepo.ReadById(deleteDto.SslPairId)
	if err != nil {
		log.Printf("SslPairNotFound: %s", err)
		return errors.New("SslPairNotFound")
	}

	err = sslCmdRepo.Delete(deleteDto.SslPairId)
	if err != nil {
		log.Printf("DeleteSslPairError: %s", err)
		return errors.New("DeleteSslPairInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		DeleteSslPair(deleteDto)

	return nil
}
