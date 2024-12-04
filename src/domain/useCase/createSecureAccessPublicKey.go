package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

func CreateSecureAccessPublicKey(
	secureAccessKeyCmdRepo repository.SecureAccessKeyCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	createDto dto.CreateSecureAccessPublicKey,
) error {
	keyId, err := secureAccessKeyCmdRepo.Create(createDto)
	if err != nil {
		slog.Error("CreateSecureAccessPublicKeyError", slog.Any("error", err))
		return errors.New("CreateSecureAccessPublicKeyInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		CreateSecureAccessPublicKey(createDto, keyId)

	return nil
}
