package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

func DeleteAccount(
	accountQueryRepo repository.AccountQueryRepo,
	accountCmdRepo repository.AccountCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	deleteDto dto.DeleteAccount,
) error {
	_, err := accountQueryRepo.ReadById(deleteDto.AccountId)
	if err != nil {
		return errors.New("AccountNotFound")
	}

	err = accountCmdRepo.Delete(deleteDto.AccountId)
	if err != nil {
		slog.Error("DeleteAccountInfraError", slog.Any("error", err))
		return errors.New("DeleteAccountInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).DeleteAccount(deleteDto)

	return nil
}
