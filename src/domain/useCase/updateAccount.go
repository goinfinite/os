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
		return err
	}

	if updateDto.Password != nil {
		err = accountCmdRepo.UpdatePassword(updateDto.AccountId, *updateDto.Password)
		if err != nil {
			slog.Error(
				"UpdateAccountPasswordError",
				slog.String("accountId", updateDto.AccountId.String()),
				slog.Any("error", err),
			)
			return errors.New("UpdateAccountPasswordInfraError")
		}

	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).UpdateAccount(updateDto)

	return nil
}
