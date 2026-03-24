package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
	tkRepository "github.com/goinfinite/tk/src/domain/repository"
)

func CreateSecureAccessPublicKey(
	accountCmdRepo repository.AccountCmdRepo,
	activityRecordCmdRepo tkRepository.ActivityRecordCmdRepo,
	createDto dto.CreateSecureAccessPublicKey,
) error {
	keyId, err := accountCmdRepo.CreateSecureAccessPublicKey(createDto)
	if err != nil {
		slog.Error("CreateSecureAccessPublicKeyError", slog.String("err", err.Error()))
		return errors.New("CreateSecureAccessPublicKeyInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		CreateSecureAccessPublicKey(createDto, keyId)

	return nil
}
