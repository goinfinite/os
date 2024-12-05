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
	readRequestDto := dto.ReadAccountsRequest{
		AccountId: &updateDto.AccountId,
	}
	_, err = accountQueryRepo.ReadFirst(readRequestDto)
	if err != nil {
		return newKey, err
	}

	newKey, err = accountCmdRepo.UpdateApiKey(updateDto.AccountId)
	if err != nil {
		slog.Error(
			"UpdateAccountApiKeyError",
			slog.String("accountId", updateDto.AccountId.String()),
			slog.Any("error", err),
		)
		return newKey, errors.New("UpdateAccountApiKeyInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).UpdateAccount(updateDto)

	return newKey, nil
}
