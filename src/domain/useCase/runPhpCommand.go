package useCase

import (
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

var RuntimeRunPhpCommandTimeoutSecsDefault uint64 = 600

func RunPhpCommand(
	runtimeCmdRepo repository.RuntimeCmdRepo,
	runRequest dto.RunPhpCommandRequest,
) (runResponse dto.RunPhpCommandResponse, err error) {
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
