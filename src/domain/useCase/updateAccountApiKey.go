package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
)

func UpdateAccountApiKey(
	accountQueryRepo repository.AccountQueryRepo,
	accountCmdRepo repository.AccountCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	updateDto dto.UpdateAccount,
) (newKey valueObject.AccessTokenStr, err error) {
	if updateDto.AccountId == nil && updateDto.AccountUsername == nil {
		return newKey, errors.New("AccountIdOrUsernameRequired")
	}

	accountEntity, err := accountQueryRepo.ReadFirst(dto.ReadAccountsRequest{
		AccountId:       updateDto.AccountId,
		AccountUsername: updateDto.AccountUsername,
	})
	if err != nil {
		return newKey, errors.New("AccountNotFound")
	}

	newKey, err = accountCmdRepo.UpdateApiKey(accountEntity.Id)
	if err != nil {
		slog.Error(
			"UpdateAccountApiKeyError",
			slog.String("accountId", accountEntity.Id.String()),
			slog.String("err", err.Error()),
		)
		return newKey, errors.New("UpdateAccountApiKeyInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).UpdateAccount(accountEntity.Id, updateDto)

	return newKey, nil
}
