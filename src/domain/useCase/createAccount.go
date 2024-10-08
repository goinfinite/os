package useCase

import (
	"errors"
	"log/slog"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
)

func CreateAccount(
	accountQueryRepo repository.AccountQueryRepo,
	accountCmdRepo repository.AccountCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	createDto dto.CreateAccount,
) error {
	_, err := accountQueryRepo.ReadByUsername(createDto.Username)
	if err == nil {
		return errors.New("AccountAlreadyExists")
	}

	_, err = accountCmdRepo.Create(createDto)
	if err != nil {
		slog.Error("CreateAccountInfraError", slog.Any("error", err))
		return errors.New("CreateAccountInfraError")
	}

	return nil
}
