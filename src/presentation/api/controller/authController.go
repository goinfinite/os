package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/sam/src/domain/dto"
	"github.com/speedianet/sam/src/domain/useCase"
	"github.com/speedianet/sam/src/domain/valueObject"
	"github.com/speedianet/sam/src/infra"
	restApiHelper "github.com/speedianet/sam/src/presentation/api/helper"
)

// AuthLogin godoc
// @Summary      Generate JWT with credentials
// @Description  Generate JWT with credentials
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        loginDto 	  body    dto.Login  true  "Login"
// @Success      200 {object} entity.AccessToken
// @Router       /auth/login/ [post]
func AuthLoginController(c echo.Context) error {
	requiredParams := []string{"username", "password"}
	requestBody, _ := restApiHelper.GetRequestBody(c)

	restApiHelper.CheckMissingParams(requestBody, requiredParams)

	loginDto := dto.NewLogin(
		valueObject.NewUsernamePanic(requestBody["username"].(string)),
		valueObject.NewPasswordPanic(requestBody["password"].(string)),
	)

	authQueryRepo := infra.AuthQueryRepo{}
	authCmdRepo := infra.AuthCmdRepo{}

	ipAddress := valueObject.NewIpAddressPanic(c.RealIP())

	accessToken := useCase.GetSessionToken(
		authQueryRepo,
		authCmdRepo,
		loginDto,
		ipAddress,
	)

	return restApiHelper.ResponseWrapper(c, http.StatusOK, accessToken)
}
