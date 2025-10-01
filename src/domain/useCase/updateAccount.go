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
	accountEntity, err := accountQueryRepo.ReadFirst(dto.ReadAccountsRequest{
		AccountId:       updateDto.AccountId,
		AccountUsername: updateDto.AccountUsername,
	})
	if err != nil {
		return errors.New("AccountNotFound")
	}

	err = accountCmdRepo.Update(updateDto)
	if err != nil {
		slog.Error(
			"UpdateAccount",
			slog.String("accountId", accountEntity.Id.String()),
			slog.String("err", err.Error()),
		)
		return errors.New("UpdateAccountInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		UpdateAccount(accountEntity.Id, updateDto)

	return nil
}
