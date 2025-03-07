package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

func CreateFirstAccount(
	accountQueryRepo repository.AccountQueryRepo,
	accountCmdRepo repository.AccountCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	createDto dto.CreateAccount,
) error {
	readRequestDto := dto.ReadAccountsRequest{}
	_, err := accountQueryRepo.ReadFirst(readRequestDto)
	if err == nil {
		return errors.New("AtLeastOneAccountAlreadyExists")
	}

	accountId, err := accountCmdRepo.Create(createDto)
	if err != nil {
		slog.Error("CreateFirstAccountError", slog.String("err", err.Error()))
		return errors.New("CreateFirstAccountInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		CreateAccount(createDto, accountId)

	return nil
}
