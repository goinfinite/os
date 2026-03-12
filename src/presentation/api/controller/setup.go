package apiController

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	accountInfra "github.com/goinfinite/os/src/infra/account"
	activityRecordInfra "github.com/goinfinite/os/src/infra/activityRecord"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	"github.com/goinfinite/os/src/presentation/liaison"
	tkPresentation "github.com/goinfinite/tk/src/presentation"
	"github.com/labstack/echo/v4"
)

type SetupController struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
	trailDbSvc      *internalDbInfra.TrailDatabaseService
}

func NewSetupController(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *SetupController {
	return &SetupController{
		persistentDbSvc: persistentDbSvc,
		trailDbSvc:      trailDbSvc,
	}
}

// SetupInfiniteOs godoc
// @Summary      SetupInfiniteOs
// @Description  Creates the first Infinite OS account without requiring authentication.<br />This can only be used when the Infinite OS interface is accessed for the first time with no accounts created.
// @Tags         setup
// @Accept       json
// @Produce      json
// @Param        createFirstAccount body dto.CreateAccount true "CreateFirstAccount"
// @Success      201 {object} object{} "FirstAccountCreated"
// @Router       /v1/setup/ [post]
func (controller *SetupController) Setup(echoContext echo.Context) error {
	requestBody, err := tkPresentation.ApiRequestInputReader{}.Reader(echoContext)
	if err != nil {
		return err
	}

	accountQueryRepo := accountInfra.NewAccountQueryRepo(controller.persistentDbSvc)
	accountCmdRepo := accountInfra.NewAccountCmdRepo(controller.persistentDbSvc)
	activityRecordCmdRepo := activityRecordInfra.NewActivityRecordCmdRepo(
		controller.trailDbSvc,
	)

	username, err := valueObject.NewUsername(requestBody["username"])
	if err != nil {
		return tkPresentation.LiaisonApiResponseEmitter(echoContext, tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error()))
	}

	password, err := tkValueObject.NewPassword(requestBody["password"])
	if err != nil {
		return tkPresentation.LiaisonApiResponseEmitter(echoContext, tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error()))
	}

	isSuperAdmin := false

	operatorIpAddress := liaison.LocalOperatorIpAddress
	if requestBody["operatorIpAddress"] != nil {
		operatorIpAddress, err = tkValueObject.NewIpAddress(
			requestBody["operatorIpAddress"],
		)
		if err != nil {
			return tkPresentation.LiaisonApiResponseEmitter(echoContext, tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error()))
		}
	}

	createDto := dto.NewCreateAccount(
		username, password, isSuperAdmin, liaison.LocalOperatorAccountId,
		operatorIpAddress,
	)

	err = useCase.CreateFirstAccount(
		accountQueryRepo, accountCmdRepo, activityRecordCmdRepo, createDto,
	)
	if err != nil {
		return tkPresentation.LiaisonApiResponseEmitter(echoContext, tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusInfraError, err.Error()))
	}

	return tkPresentation.LiaisonApiResponseEmitter(echoContext, tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusCreated, "FirstAccountCreated"))
}
