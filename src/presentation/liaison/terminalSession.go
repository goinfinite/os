package liaison

import (
	"errors"

	tkPresentation "github.com/goinfinite/tk/src/presentation"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	accountInfra "github.com/goinfinite/os/src/infra/account"
	activityRecordInfra "github.com/goinfinite/os/src/infra/activityRecord"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	terminalSessionInfra "github.com/goinfinite/os/src/infra/terminalSession"
	liaisonHelper "github.com/goinfinite/os/src/presentation/liaison/helper"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type TerminalSessionLiaison struct {
	accountQueryRepo         *accountInfra.AccountQueryRepo
	terminalSessionQueryRepo *terminalSessionInfra.TerminalSessionQueryRepo
	terminalSessionCmdRepo   *terminalSessionInfra.TerminalSessionCmdRepo
	activityRecordCmdRepo    *activityRecordInfra.ActivityRecordCmdRepo
}

func NewTerminalSessionLiaison(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *TerminalSessionLiaison {
	accountQueryRepo := accountInfra.NewAccountQueryRepo(persistentDbSvc)

	return &TerminalSessionLiaison{
		accountQueryRepo:         accountQueryRepo,
		terminalSessionQueryRepo: terminalSessionInfra.NewTerminalSessionQueryRepo(accountQueryRepo),
		terminalSessionCmdRepo:   terminalSessionInfra.NewTerminalSessionCmdRepo(),
		activityRecordCmdRepo:    activityRecordInfra.NewActivityRecordCmdRepo(trailDbSvc),
	}
}

func (liaison *TerminalSessionLiaison) readTerminalSessionById(
	terminalSessionId valueObject.TerminalSessionId,
	operatorAccountId tkValueObject.AccountId,
) (terminalSessionEntity entity.TerminalSession, err error) {
	return useCase.NewReadTerminalSessions(
		liaison.accountQueryRepo, liaison.terminalSessionQueryRepo,
	).ExecuteFirst(dto.ReadTerminalSessionsRequest{
		TerminalSessionId: &terminalSessionId,
		OperatorAccountId: operatorAccountId,
	})
}

func (liaison *TerminalSessionLiaison) Read(
	untrustedInput map[string]any,
) tkPresentation.LiaisonResponse {
	operatorAccountId, _, err := liaisonHelper.ReadOperatorContext(untrustedInput)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError, err.Error(),
		)
	}

	var terminalSessionIdPtr *valueObject.TerminalSessionId
	if untrustedInput["id"] != nil {
		terminalSessionId, err := valueObject.NewTerminalSessionId(untrustedInput["id"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError, err.Error(),
			)
		}
		terminalSessionIdPtr = &terminalSessionId
	}

	var accountIdPtr *tkValueObject.AccountId
	if untrustedInput["accountId"] != nil {
		accountId, err := tkValueObject.NewAccountId(untrustedInput["accountId"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError, err.Error(),
			)
		}
		accountIdPtr = &accountId
	}

	requestPagination, err := tkPresentation.PaginationParser(
		useCase.TerminalSessionsDefaultPagination, untrustedInput,
	)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError, err.Error(),
		)
	}

	responseDto, err := useCase.NewReadTerminalSessions(
		liaison.accountQueryRepo, liaison.terminalSessionQueryRepo,
	).Execute(dto.ReadTerminalSessionsRequest{
		Pagination:        requestPagination,
		TerminalSessionId: terminalSessionIdPtr,
		AccountId:         accountIdPtr,
		OperatorAccountId: operatorAccountId,
	})
	if err != nil {
		if errors.Is(err, useCase.ErrTerminalSessionAccountNotOwned) {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError, err.Error(),
			)
		}

		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusInfraError, err.Error(),
		)
	}

	return tkPresentation.NewLiaisonResponseNoMessage(
		tkPresentation.LiaisonResponseStatusSuccess, responseDto,
	)
}

func (liaison *TerminalSessionLiaison) ReadFirst(
	untrustedInput map[string]any,
) tkPresentation.LiaisonResponse {
	operatorAccountId, _, err := liaisonHelper.ReadOperatorContext(untrustedInput)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError, err.Error(),
		)
	}

	terminalSessionId, err := valueObject.NewTerminalSessionId(untrustedInput["id"])
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError, err.Error(),
		)
	}

	terminalSessionEntity, err := liaison.readTerminalSessionById(
		terminalSessionId, operatorAccountId,
	)
	if err != nil {
		if errors.Is(err, repository.ErrTerminalSessionNotFound) {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError, err.Error(),
			)
		}

		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusInfraError, err.Error(),
		)
	}

	return tkPresentation.NewLiaisonResponseNoMessage(
		tkPresentation.LiaisonResponseStatusSuccess, terminalSessionEntity,
	)
}

func (liaison *TerminalSessionLiaison) Create(
	untrustedInput map[string]any,
) tkPresentation.LiaisonResponse {
	operatorAccountId, operatorIpAddress, err := liaisonHelper.ReadOperatorContext(
		untrustedInput,
	)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError, err.Error(),
		)
	}

	var workingDirPtr *tkValueObject.UnixAbsoluteFilePath
	if untrustedInput["workingDir"] != nil && untrustedInput["workingDir"] != "" {
		workingDir, err := tkValueObject.NewUnixAbsoluteFilePath(
			untrustedInput["workingDir"], false,
		)
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError, err.Error(),
			)
		}
		workingDirPtr = &workingDir
	}

	var commandPtr *tkValueObject.UnixCommand
	if untrustedInput["command"] != nil && untrustedInput["command"] != "" {
		command, err := tkValueObject.NewUnixCommand(untrustedInput["command"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError, err.Error(),
			)
		}
		commandPtr = &command
	}

	var accountIdPtr *tkValueObject.AccountId
	if untrustedInput["accountId"] != nil && untrustedInput["accountId"] != "" {
		accountId, err := tkValueObject.NewAccountId(untrustedInput["accountId"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError, err.Error(),
			)
		}
		accountIdPtr = &accountId
	}

	createDto := dto.NewCreateTerminalSession(
		workingDirPtr, commandPtr, accountIdPtr, operatorAccountId, operatorIpAddress,
	)

	terminalSessionId, err := useCase.NewCreateTerminalSession(
		liaison.accountQueryRepo, liaison.terminalSessionQueryRepo,
		liaison.terminalSessionCmdRepo, liaison.activityRecordCmdRepo,
		infraHelper.ReadMaxTerminalSessionsPerAccount(),
	).Execute(createDto)
	if err != nil {
		if liaisonHelper.IsAnyError(
			err,
			repository.ErrTerminalSessionAccountCapReached,
			repository.ErrTerminalSessionWorkingDirNotFound,
			repository.ErrTerminalSessionAccountRequired,
		) {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError, err.Error(),
			)
		}

		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusInfraError, err.Error(),
		)
	}

	terminalSessionEntity, err := liaison.readTerminalSessionById(
		terminalSessionId, operatorAccountId,
	)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusInfraError, err.Error(),
		)
	}

	return tkPresentation.NewLiaisonResponseNoMessage(
		tkPresentation.LiaisonResponseStatusCreated, terminalSessionEntity,
	)
}

func (liaison *TerminalSessionLiaison) Delete(
	untrustedInput map[string]any,
) tkPresentation.LiaisonResponse {
	operatorAccountId, operatorIpAddress, err := liaisonHelper.ReadOperatorContext(
		untrustedInput,
	)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError, err.Error(),
		)
	}

	terminalSessionId, err := valueObject.NewTerminalSessionId(untrustedInput["id"])
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError, err.Error(),
		)
	}

	deleteDto := dto.NewDeleteTerminalSession(
		terminalSessionId, operatorAccountId, operatorIpAddress,
	)

	err = useCase.DeleteTerminalSession(
		liaison.accountQueryRepo, liaison.terminalSessionQueryRepo,
		liaison.terminalSessionCmdRepo, liaison.activityRecordCmdRepo, deleteDto,
	)
	if err != nil {
		if errors.Is(err, repository.ErrTerminalSessionNotFound) {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError, err.Error(),
			)
		}

		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusInfraError, err.Error(),
		)
	}

	return tkPresentation.NewLiaisonResponseNoMessage(
		tkPresentation.LiaisonResponseStatusSuccess, "TerminalSessionDeleted",
	)
}

func (liaison *TerminalSessionLiaison) Attach(untrustedInput map[string]any) (
	attachHandle repository.TerminalSessionAttachHandle,
	response tkPresentation.LiaisonResponse,
) {
	operatorAccountId, _, err := liaisonHelper.ReadOperatorContext(untrustedInput)
	if err != nil {
		return attachHandle, tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError, err.Error(),
		)
	}

	terminalSessionId, err := valueObject.NewTerminalSessionId(untrustedInput["id"])
	if err != nil {
		return attachHandle, tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError, err.Error(),
		)
	}

	attachDto := dto.NewAttachTerminalSession(terminalSessionId, operatorAccountId)

	attachHandle, err = useCase.NewAttachTerminalSession(
		liaison.accountQueryRepo, liaison.terminalSessionQueryRepo,
		liaison.terminalSessionCmdRepo, infraHelper.IsReadOnlyMode(),
	).Execute(attachDto)
	if err != nil {
		if liaisonHelper.IsAnyError(
			err,
			useCase.ErrReadOnlyMode,
			repository.ErrTerminalSessionNotFound,
			useCase.ErrTerminalSessionAccountNotOwned,
		) {
			return attachHandle, tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError, err.Error(),
			)
		}

		return attachHandle, tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusInfraError, err.Error(),
		)
	}

	return attachHandle, tkPresentation.NewLiaisonResponseNoMessage(
		tkPresentation.LiaisonResponseStatusSuccess, "TerminalSessionAttached",
	)
}
