package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

func CreateSecureAccessKey(
	accountQueryRepo repository.AccountQueryRepo,
	secureAccessKeyCmdRepo repository.SecureAccessKeyCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	createDto dto.CreateSecureAccessKey,
) error {
	_, err := accountQueryRepo.ReadById(createDto.AccountId)
	if err != nil {
		return errors.New("AccountNotFound")
	}

	keyId, err := secureAccessKeyCmdRepo.Create(createDto)
	if err != nil {
		slog.Error("CreateSecureAccessKeyError", slog.Any("error", err))
		return errors.New("CreateSecureAccessKeyInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		CreateSecureAccessKey(createDto, keyId)

	return nil
}
