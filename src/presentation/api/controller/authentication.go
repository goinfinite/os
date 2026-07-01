package apiController

import (
	_ "github.com/goinfinite/os/src/domain/dto"
	_ "github.com/goinfinite/os/src/domain/entity"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	"github.com/goinfinite/os/src/presentation/liaison"
	tkPresentation "github.com/goinfinite/tk/src/presentation"
	"github.com/labstack/echo/v4"
)

type AuthenticationController struct {
	authenticationLiaison *liaison.AuthenticationLiaison
}

func NewAuthenticationController(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *AuthenticationController {
	return &AuthenticationController{
		authenticationLiaison: liaison.NewAuthenticationLiaison(
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
	inputReader := tkPresentation.ApiRequestInputReader{}
	requestData, requestParsingErr := inputReader.Reader(echoContext)
	if requestParsingErr != nil {
		return requestParsingErr
	}

	return tkPresentation.LiaisonApiResponseEmitter(
		echoContext, controller.authenticationLiaison.Login(requestData),
	)
}
