package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/repository"
	useCaseHelper "github.com/goinfinite/os/src/domain/useCase/helper"
	tkDto "github.com/goinfinite/tk/src/domain/dto"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

var (
	ErrTerminalSessionAccountNotOwned = errors.New("TerminalSessionAccountNotOwned")

	TerminalSessionsDefaultSortBy tkValueObject.PaginationSortBy = tkValueObject.PaginationSortBy("createdAt")

	terminalSessionsDefaultSortDirection tkValueObject.PaginationSortDirection = tkValueObject.PaginationSortDirection("desc")

	TerminalSessionsDefaultPagination tkDto.Pagination = tkDto.Pagination{
		PageNumber:    0,
		ItemsPerPage:  10,
		SortBy:        &TerminalSessionsDefaultSortBy,
		SortDirection: &terminalSessionsDefaultSortDirection,
	}
)

type ReadTerminalSessions struct {
	accountQueryRepo         repository.AccountQueryRepo
	terminalSessionQueryRepo repository.TerminalSessionQueryRepo
}

func NewReadTerminalSessions(
	accountQueryRepo repository.AccountQueryRepo,
	terminalSessionQueryRepo repository.TerminalSessionQueryRepo,
) ReadTerminalSessions {
	return ReadTerminalSessions{
		accountQueryRepo:         accountQueryRepo,
		terminalSessionQueryRepo: terminalSessionQueryRepo,
	}
}

func (uc ReadTerminalSessions) authorizeRead(
	requestDto dto.ReadTerminalSessionsRequest,
) (authorizedRequestDto dto.ReadTerminalSessionsRequest, err error) {
	operatorAccountEntity, isSystemOperator, err := useCaseHelper.ReadOperatorAccount(
		uc.accountQueryRepo, requestDto.OperatorAccountId,
	)
	if err != nil {
		slog.Error("ReadOperatorAccountError", slog.String("err", err.Error()))
		return requestDto, errors.New("ReadOperatorAccountInfraError")
	}

	if isSystemOperator || operatorAccountEntity.IsSuperAdmin {
		return requestDto, nil
	}

	if requestDto.AccountId != nil && *requestDto.AccountId != operatorAccountEntity.Id {
		return requestDto, ErrTerminalSessionAccountNotOwned
	}

	operatorAccountId := operatorAccountEntity.Id
	requestDto.AccountId = &operatorAccountId

	return requestDto, nil
}

func (uc ReadTerminalSessions) Execute(
	requestDto dto.ReadTerminalSessionsRequest,
) (responseDto dto.ReadTerminalSessionsResponse, err error) {
	authorizedRequestDto, err := uc.authorizeRead(requestDto)
	if err != nil {
		return responseDto, err
	}

	responseDto, err = uc.terminalSessionQueryRepo.Read(authorizedRequestDto)
	if err != nil {
		slog.Error("ReadTerminalSessionsError", slog.String("err", err.Error()))
		return responseDto, errors.New("ReadTerminalSessionsInfraError")
	}

	return responseDto, nil
}

func (uc ReadTerminalSessions) ExecuteFirst(
	requestDto dto.ReadTerminalSessionsRequest,
) (terminalSessionEntity entity.TerminalSession, err error) {
	authorizedRequestDto, err := uc.authorizeRead(requestDto)
	if err != nil {
		return terminalSessionEntity, err
	}

	return uc.terminalSessionQueryRepo.ReadFirst(authorizedRequestDto)
}
