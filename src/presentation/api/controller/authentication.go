package apiController

import (
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	apiHelper "github.com/goinfinite/os/src/presentation/api/helper"
	"github.com/goinfinite/os/src/presentation/service"
	"github.com/labstack/echo/v4"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *AuthController {
	return &AuthController{
		authService: service.NewAuthService(trailDbSvc),
	}
}

// GenerateJwtWithCredentials godoc
// @Summary      GenerateJwtWithCredentials
// @Description  Generate JWT with credentials
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        loginDto 	  body    dto.Login  true  "All props are required."
// @Success      200 {object} entity.AccessToken
// @Failure      401 {object} string
// @Router       /v1/auth/login/ [post]
func (controller *AuthController) GenerateJwtWithCredentials(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}
	requestBody["ipAddress"] = c.RealIP()

	return apiHelper.ServiceResponseWrapper(
		c, controller.authService.GenerateJwtWithCredentials(requestBody),
	)
}
