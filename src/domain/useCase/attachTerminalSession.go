package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

var ErrReadOnlyMode = errors.New("ReadOnlyModeEnabled")

type AttachTerminalSession struct {
	accountQueryRepo         repository.AccountQueryRepo
	terminalSessionQueryRepo repository.TerminalSessionQueryRepo
	terminalSessionCmdRepo   repository.TerminalSessionCmdRepo
	isReadOnlyMode           bool
}

func NewAttachTerminalSession(
	accountQueryRepo repository.AccountQueryRepo,
	terminalSessionQueryRepo repository.TerminalSessionQueryRepo,
	terminalSessionCmdRepo repository.TerminalSessionCmdRepo,
	isReadOnlyMode bool,
) AttachTerminalSession {
	return AttachTerminalSession{
		accountQueryRepo:         accountQueryRepo,
		terminalSessionQueryRepo: terminalSessionQueryRepo,
		terminalSessionCmdRepo:   terminalSessionCmdRepo,
		isReadOnlyMode:           isReadOnlyMode,
	}
}

func (uc AttachTerminalSession) Execute(
	attachDto dto.AttachTerminalSession,
) (attachHandle repository.TerminalSessionAttachHandle, err error) {
	if uc.isReadOnlyMode {
		return attachHandle, ErrReadOnlyMode
	}

	terminalSessionEntity, err := NewReadTerminalSessions(
		uc.accountQueryRepo, uc.terminalSessionQueryRepo,
	).ExecuteFirst(dto.ReadTerminalSessionsRequest{
		TerminalSessionId: &attachDto.Id,
		OperatorAccountId: attachDto.OperatorAccountId,
	})
	if err != nil {
		return attachHandle, err
	}

	attachDto.AccountUsername = terminalSessionEntity.AccountUsername

	attachHandle, err = uc.terminalSessionCmdRepo.Attach(attachDto)
	if err != nil {
		slog.Error("AttachTerminalSessionError", slog.String("err", err.Error()))
		return attachHandle, errors.New("AttachTerminalSessionInfraError")
	}

	return attachHandle, nil
}
