package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

func DeleteSecureAccessKey(
	secureAccessKeyQueryRepo repository.SecureAccessKeyQueryRepo,
	secureAccessKeyCmdRepo repository.SecureAccessKeyCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	deleteDto dto.DeleteSecureAccessKey,
) error {
	readRequestDto := dto.ReadSecureAccessKeysRequest{
		SecureAccessKeyId: &deleteDto.Id,
	}
	keyToDelete, err := secureAccessKeyQueryRepo.ReadFirst(readRequestDto)
	if err != nil {
		return errors.New("SecureAccessKeyNotFound")
	}

	err = secureAccessKeyCmdRepo.Delete(keyToDelete.Id)
	if err != nil {
		slog.Error("DeleteSecureAccessKeyError", slog.Any("error", err))
		return errors.New("DeleteSecureAccessKeyInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		DeleteSecureAccessKey(deleteDto, keyToDelete.AccountId)

	return nil
}
