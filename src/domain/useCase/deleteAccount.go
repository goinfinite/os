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

	readResponseDto, err := accountQueryRepo.Read(
		dto.ReadAccountsRequest{Pagination: AccountsDefaultPagination},
	)
	if err != nil {
		slog.Error("ReadAccountsError", slog.String("err", err.Error()))
		return errors.New("ReadAccountsInfraError")
	}

	if len(readResponseDto.Accounts) <= 1 {
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
