package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

func UpdateAccount(
	accountQueryRepo repository.AccountQueryRepo,
	accountCmdRepo repository.AccountCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	updateDto dto.UpdateAccount,
) error {
	readRequestDto := dto.ReadAccountsRequest{
		AccountId: &updateDto.AccountId,
	}
	_, err := accountQueryRepo.ReadFirst(readRequestDto)
	if err != nil {
		return errors.New("AccountNotFound")
	}

	err = accountCmdRepo.Update(updateDto)
	if err != nil {
		slog.Error(
			"UpdateAccount",
			slog.String("accountId", updateDto.AccountId.String()),
			slog.String("err", err.Error()),
		)
		return errors.New("UpdateAccountInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).UpdateAccount(updateDto)

	return nil
}
