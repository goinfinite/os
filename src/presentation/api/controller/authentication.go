package apiController

import (
	"net/http"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	accountInfra "github.com/goinfinite/os/src/infra/account"
	activityRecordInfra "github.com/goinfinite/os/src/infra/activityRecord"
	authInfra "github.com/goinfinite/os/src/infra/auth"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	apiHelper "github.com/goinfinite/os/src/presentation/api/helper"
	serviceHelper "github.com/goinfinite/os/src/presentation/service/helper"
	"github.com/labstack/echo/v4"
)

type AuthController struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
	trailDbSvc      *internalDbInfra.TrailDatabaseService
}

func NewAuthController(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *AuthController {
	return &AuthController{
		persistentDbSvc: persistentDbSvc,
		trailDbSvc:      trailDbSvc,
	}
}

// CreateSessionTokenWithCredentials godoc
// @Summary      CreateSessionTokenWithCredentials
// @Description  Create a new session token with the provided credentials.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        createSessionToken body dto.CreateSessionToken true "CreateSessionToken"
// @Success      200 {object} entity.AccessToken
// @Failure      401 {object} string
// @Router       /v1/auth/login/ [post]
func (controller *AuthController) Login(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	requiredParams := []string{"username", "password"}
	err = serviceHelper.RequiredParamsInspector(requestBody, requiredParams)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	username, err := valueObject.NewUsername(requestBody["username"])
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	password, err := valueObject.NewPassword(requestBody["password"])
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	operatorIpAddress, err := valueObject.NewIpAddress(requestBody["operatorIpAddress"])
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	dto := dto.NewCreateSessionToken(username, password, operatorIpAddress)

	authQueryRepo := authInfra.NewAuthQueryRepo(controller.persistentDbSvc)
	authCmdRepo := authInfra.AuthCmdRepo{}
	accountQueryRepo := accountInfra.NewAccountQueryRepo(controller.persistentDbSvc)
	activityRecordQueryRepo := activityRecordInfra.NewActivityRecordQueryRepo(
		controller.trailDbSvc,
	)
	activityRecordCmdRepo := activityRecordInfra.NewActivityRecordCmdRepo(
		controller.trailDbSvc,
	)

	accessToken, err := useCase.CreateSessionToken(
		authQueryRepo, authCmdRepo, accountQueryRepo, activityRecordQueryRepo,
		activityRecordCmdRepo, dto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusUnauthorized, err.Error())
	}

	c.Response().Header().Set("X-Ignore-Toast", "true")
	return apiHelper.ResponseWrapper(c, http.StatusOK, accessToken)
}
