package apiController

import (
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	apiHelper "github.com/goinfinite/os/src/presentation/api/helper"
	"github.com/goinfinite/os/src/presentation/service"
	tkPresentation "github.com/goinfinite/tk/src/presentation"
	"github.com/labstack/echo/v4"
)

type AuthenticationController struct {
	authenticationService *service.AuthenticationService
}

func NewAuthenticationController(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *AuthenticationController {
	return &AuthenticationController{
		authenticationService: service.NewAuthenticationService(
			persistentDbSvc, trailDbSvc,
		),
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
func (controller *AuthenticationController) Login(echoContext echo.Context) error {
	return apiHelper.ServiceResponseWrapper(
		echoContext, controller.authenticationService.Login(
			echoContext.Get("RequestInputParsed").(tkPresentation.RequestInputParsed),
		),
	)
}
