package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

func CreateAccount(
	accountQueryRepo repository.AccountQueryRepo,
	accountCmdRepo repository.AccountCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	createDto dto.CreateAccount,
) error {
	readRequestDto := dto.ReadAccountsRequest{
		AccountUsername: &createDto.Username,
	}
	_, err := accountQueryRepo.ReadFirst(readRequestDto)
	if err == nil {
		return errors.New("AccountAlreadyExists")
	}

	accountId, err := accountCmdRepo.Create(createDto)
	if err != nil {
		slog.Error("CreateAccountError", slog.String("err", err.Error()))
		return errors.New("CreateAccountInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		CreateAccount(createDto, accountId)

	return nil
}
