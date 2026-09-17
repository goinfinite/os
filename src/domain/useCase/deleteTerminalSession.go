package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
	tkRepository "github.com/goinfinite/tk/src/domain/repository"
)

func DeleteTerminalSession(
	accountQueryRepo repository.AccountQueryRepo,
	terminalSessionQueryRepo repository.TerminalSessionQueryRepo,
	terminalSessionCmdRepo repository.TerminalSessionCmdRepo,
	activityRecordCmdRepo tkRepository.ActivityRecordCmdRepo,
	deleteDto dto.DeleteTerminalSession,
) error {
	terminalSessionEntity, err := NewReadTerminalSessions(
		accountQueryRepo, terminalSessionQueryRepo,
	).ExecuteFirst(dto.ReadTerminalSessionsRequest{
		TerminalSessionId: &deleteDto.Id,
		OperatorAccountId: deleteDto.OperatorAccountId,
	})
	if err != nil {
		return err
	}

	deleteDto.AccountUsername = terminalSessionEntity.AccountUsername

	err = terminalSessionCmdRepo.Delete(deleteDto)
	if err != nil {
		slog.Error("DeleteTerminalSessionError", slog.String("err", err.Error()))
		return errors.New("DeleteTerminalSessionInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		DeleteTerminalSession(deleteDto, terminalSessionEntity.AccountId)

	return nil
}
