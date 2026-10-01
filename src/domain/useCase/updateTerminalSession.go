package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
	tkRepository "github.com/goinfinite/tk/src/domain/repository"
)

func UpdateTerminalSession(
	accountQueryRepo repository.AccountQueryRepo,
	terminalSessionQueryRepo repository.TerminalSessionQueryRepo,
	terminalSessionCmdRepo repository.TerminalSessionCmdRepo,
	activityRecordCmdRepo tkRepository.ActivityRecordCmdRepo,
	updateDto dto.UpdateTerminalSession,
) error {
	terminalSessionEntity, err := NewReadTerminalSessions(
		accountQueryRepo, terminalSessionQueryRepo,
	).ExecuteFirst(dto.ReadTerminalSessionsRequest{
		TerminalSessionId: &updateDto.Id,
		OperatorAccountId: updateDto.OperatorAccountId,
	})
	if err != nil {
		return err
	}

	updateDto.AccountUsername = terminalSessionEntity.AccountUsername

	err = terminalSessionCmdRepo.Update(updateDto)
	if err != nil {
		slog.Error("UpdateTerminalSessionError", slog.String("err", err.Error()))
		return errors.New("UpdateTerminalSessionInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		UpdateTerminalSession(updateDto, terminalSessionEntity.AccountId)

	return nil
}
