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
	_, err := accountQueryRepo.ReadFirst(
		dto.ReadAccountsRequest{AccountId: &deleteDto.AccountId},
	)
	if err != nil {
		return errors.New("AccountNotFound")
	}

	accountsCount, err := accountQueryRepo.Count(dto.ReadAccountsRequest{})
	if err != nil {
		slog.Error("CountAccountsError", slog.String("err", err.Error()))
		return errors.New("CountAccountsInfraError")
	}

	if accountsCount <= 1 {
		return errors.New("AtLeastOneAccountMustExist")
	}

	err = accountCmdRepo.Delete(deleteDto.AccountId)
	if err != nil {
		slog.Error("DeleteAccountError", slog.String("err", err.Error()))
		return errors.New("DeleteAccountInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).DeleteAccount(deleteDto)

	return nil
}
