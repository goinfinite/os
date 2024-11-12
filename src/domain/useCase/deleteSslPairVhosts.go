package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

func DeleteSslPairVhosts(
	sslQueryRepo repository.SslQueryRepo,
	sslCmdRepo repository.SslCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	deleteDto dto.DeleteSslPairVhosts,
) error {
	_, err := sslQueryRepo.ReadById(deleteDto.SslPairId)
	if err != nil {
		return errors.New("SslPairNotFound")
	}

	err = sslCmdRepo.DeleteSslPairVhosts(deleteDto)
	if err != nil {
		slog.Error("DeleteSslPairVhostsError", slog.Any("err", err))
		return errors.New("DeleteSslPairVhostsInfraError")
	}

	return nil
}
