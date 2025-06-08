package apiController

import (
	"net/http"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	accountInfra "github.com/goinfinite/os/src/infra/account"
	activityRecordInfra "github.com/goinfinite/os/src/infra/activityRecord"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	apiHelper "github.com/goinfinite/os/src/presentation/api/helper"
	"github.com/goinfinite/os/src/presentation/liaison"
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
func (controller *SetupController) Setup(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestInputData(c)
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
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	password, err := valueObject.NewPassword(requestBody["password"])
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	isSuperAdmin := false

	operatorIpAddress := liaison.LocalOperatorIpAddress
	if requestBody["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(
			requestBody["operatorIpAddress"],
		)
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
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
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusCreated, "FirstAccountCreated")
}
