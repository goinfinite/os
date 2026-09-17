package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/repository"
	useCaseHelper "github.com/goinfinite/os/src/domain/useCase/helper"
	"github.com/goinfinite/os/src/domain/valueObject"
	tkDto "github.com/goinfinite/tk/src/domain/dto"
	tkRepository "github.com/goinfinite/tk/src/domain/repository"
)

type CreateTerminalSession struct {
	accountQueryRepo              repository.AccountQueryRepo
	terminalSessionQueryRepo      repository.TerminalSessionQueryRepo
	terminalSessionCmdRepo        repository.TerminalSessionCmdRepo
	activityRecordCmdRepo         tkRepository.ActivityRecordCmdRepo
	maxTerminalSessionsPerAccount uint16
}

func NewCreateTerminalSession(
	accountQueryRepo repository.AccountQueryRepo,
	terminalSessionQueryRepo repository.TerminalSessionQueryRepo,
	terminalSessionCmdRepo repository.TerminalSessionCmdRepo,
	activityRecordCmdRepo tkRepository.ActivityRecordCmdRepo,
	maxTerminalSessionsPerAccount uint16,
) CreateTerminalSession {
	return CreateTerminalSession{
		accountQueryRepo:              accountQueryRepo,
		terminalSessionQueryRepo:      terminalSessionQueryRepo,
		terminalSessionCmdRepo:        terminalSessionCmdRepo,
		activityRecordCmdRepo:         activityRecordCmdRepo,
		maxTerminalSessionsPerAccount: maxTerminalSessionsPerAccount,
	}
}

func (uc CreateTerminalSession) resolveOwnerAccount(
	createDto dto.CreateTerminalSession,
) (ownerAccountEntity entity.Account, err error) {
	operatorAccountEntity, isSystemOperator, err := useCaseHelper.ReadOperatorAccount(
		uc.accountQueryRepo, createDto.OperatorAccountId,
	)
	if err != nil {
		return ownerAccountEntity, errors.New("ReadOperatorAccountInfraError")
	}

	if !isSystemOperator && !operatorAccountEntity.IsSuperAdmin {
		return operatorAccountEntity, nil
	}

	if createDto.AccountId == nil {
		if isSystemOperator {
			return ownerAccountEntity, repository.ErrTerminalSessionAccountRequired
		}

		return operatorAccountEntity, nil
	}

	ownerAccountEntity, err = uc.accountQueryRepo.ReadFirst(dto.ReadAccountsRequest{
		AccountId: createDto.AccountId,
	})
	if err != nil {
		return ownerAccountEntity, errors.New("ReadOwnerAccountInfraError")
	}

	return ownerAccountEntity, nil
}

func (uc CreateTerminalSession) Execute(
	createDto dto.CreateTerminalSession,
) (terminalSessionId valueObject.TerminalSessionId, err error) {
	if createDto.WorkingDir == nil {
		defaultWorkingDir := valueObject.UnixFilePathAppWorkingDir
		createDto.WorkingDir = &defaultWorkingDir
	}

	ownerAccountEntity, err := uc.resolveOwnerAccount(createDto)
	if err != nil {
		slog.Error("ResolveTerminalSessionOwnerError", slog.String("err", err.Error()))
		return terminalSessionId, err
	}

	createDto.AccountUsername = ownerAccountEntity.Username

	ownerTerminalSessionsResponse, err := uc.terminalSessionQueryRepo.Read(
		dto.ReadTerminalSessionsRequest{
			Pagination: tkDto.PaginationUnpaginated,
			AccountId:  &ownerAccountEntity.Id,
		},
	)
	if err != nil {
		slog.Error("ReadTerminalSessionsError", slog.String("err", err.Error()))
		return terminalSessionId, errors.New("ReadTerminalSessionsInfraError")
	}
	if len(ownerTerminalSessionsResponse.TerminalSessions) >= int(uc.maxTerminalSessionsPerAccount) {
		return terminalSessionId, repository.ErrTerminalSessionAccountCapReached
	}

	terminalSessionId, err = uc.terminalSessionCmdRepo.Create(createDto)
	if err != nil {
		slog.Error("CreateTerminalSessionError", slog.String("err", err.Error()))
		if errors.Is(err, repository.ErrTerminalSessionWorkingDirNotFound) {
			return terminalSessionId, err
		}
		return terminalSessionId, errors.New("CreateTerminalSessionInfraError")
	}

	NewCreateSecurityActivityRecord(uc.activityRecordCmdRepo).
		CreateTerminalSession(createDto, ownerAccountEntity.Id, terminalSessionId)

	return terminalSessionId, nil
}
