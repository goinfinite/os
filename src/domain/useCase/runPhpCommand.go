package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
)

var RuntimeRunPhpCommandTimeoutSecsDefault uint64 = 600

func RunPhpCommand(
	accountQueryRepo repository.AccountQueryRepo,
	runtimeCmdRepo repository.RuntimeCmdRepo,
	runRequest dto.RunPhpCommandRequest,
) (runResponse dto.RunPhpCommandResponse, err error) {
	if runRequest.OperatorAccountId != valueObject.AccountIdSystem {
		operatorAccountEntity, err := accountQueryRepo.ReadFirst(
			dto.ReadAccountsRequest{AccountId: &runRequest.OperatorAccountId},
		)
		if err != nil {
			return runResponse, err
		}

		if !operatorAccountEntity.IsSuperAdmin {
			return runResponse, errors.New("OperatorIsNotSuperAdmin")
		}
	}

	if runRequest.TimeoutSecs == nil {
		runRequest.TimeoutSecs = &RuntimeRunPhpCommandTimeoutSecsDefault
	}

	runResponse, err = runtimeCmdRepo.RunPhpCommand(runRequest)
	if err != nil {
		slog.Error("RunPhpCommandError", slog.String("err", err.Error()))
		return runResponse, err
	}

	return runResponse, nil
}
