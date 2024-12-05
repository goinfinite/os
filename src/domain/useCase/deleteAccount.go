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
	readRequestDto := dto.ReadAccountsRequest{
		AccountId: &deleteDto.AccountId,
	}
	_, err := accountQueryRepo.ReadFirst(readRequestDto)
	if err != nil {
		return err
	}

	err = accountCmdRepo.Delete(deleteDto.AccountId)
	if err != nil {
		slog.Error("DeleteAccountError", slog.Any("error", err))
		return errors.New("DeleteAccountInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).DeleteAccount(deleteDto)

	return nil
}
