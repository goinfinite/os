package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

func DeleteSslPair(
	sslQueryRepo repository.SslQueryRepo,
	sslCmdRepo repository.SslCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	deleteDto dto.DeleteSslPair,
) error {
	_, err := sslQueryRepo.ReadFirst(dto.ReadSslPairsRequest{
		SslPairId: &deleteDto.SslPairId,
	})
	if err != nil {
		slog.Error("ReadSslPairError", slog.String("error", err.Error()))
		return errors.New("SslPairNotFound")
	}

	err = sslCmdRepo.Delete(deleteDto.SslPairId)
	if err != nil {
		slog.Error("DeleteSslPairError", slog.String("error", err.Error()))
		return errors.New("DeleteSslPairInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		DeleteSslPair(deleteDto)

	return nil
}
