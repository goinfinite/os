package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

func DeleteSecureAccessKey(
	accountQueryRepo repository.AccountQueryRepo,
	secureAccessKeyCmdRepo repository.SecureAccessKeyCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	deleteDto dto.DeleteSecureAccessKey,
) error {
	_, err := accountQueryRepo.ReadById(deleteDto.AccountId)
	if err != nil {
		return errors.New("AccountNotFound")
	}

	err = secureAccessKeyCmdRepo.Delete(deleteDto)
	if err != nil {
		slog.Error("DeleteSecureAccessKeyError", slog.Any("error", err))
		return errors.New("DeleteSecureAccessKeyInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		DeleteSecureAccessKey(deleteDto)

	return nil
}
