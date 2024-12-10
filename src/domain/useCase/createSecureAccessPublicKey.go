package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

func CreateSecureAccessPublicKey(
	accountCmdRepo repository.AccountCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	createDto dto.CreateSecureAccessPublicKey,
) error {
	keyId, err := accountCmdRepo.CreateSecureAccessPublicKey(createDto)
	if err != nil {
		slog.Error("CreateSecureAccessPublicKeyError", slog.Any("error", err))
		return errors.New("CreateSecureAccessPublicKeyInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		CreateSecureAccessPublicKey(createDto, keyId)

	return nil
}
