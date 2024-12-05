package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

func DeleteSecureAccessPublicKey(
	accountQueryRepo repository.AccountQueryRepo,
	accountCmdRepo repository.AccountCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	deleteDto dto.DeleteSecureAccessPublicKey,
) error {
	readRequestDto := dto.ReadSecureAccessPublicKeysRequest{
		SecureAccessPublicKeyId: &deleteDto.Id,
	}
	keyToDelete, err := accountQueryRepo.ReadFirstSecureAccessPublicKey(
		readRequestDto,
	)
	if err != nil {
		return errors.New("SecureAccessPublicKeyNotFound")
	}

	err = accountCmdRepo.DeleteSecureAccessPublicKey(keyToDelete.Id)
	if err != nil {
		slog.Error("DeleteSecureAccessPublicKeyError", slog.Any("error", err))
		return errors.New("DeleteSecureAccessPublicKeyInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		DeleteSecureAccessPublicKey(deleteDto, keyToDelete.AccountId)

	return nil
}
