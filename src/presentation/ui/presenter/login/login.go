package uiPresenter

import (
	"net/http"

	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	"github.com/goinfinite/os/src/presentation/service"
	layoutLogin "github.com/goinfinite/os/src/presentation/ui/layout/login"
	presenterHelper "github.com/goinfinite/os/src/presentation/ui/presenter/helper"
	"github.com/labstack/echo/v4"
)

type LoginPresenter struct {
	accountService *service.AccountService
}

func NewLoginPresenter(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *LoginPresenter {
	return &LoginPresenter{
		accountService: service.NewAccountService(persistentDbSvc, trailDbSvc),
	}
}

func (presenter *LoginPresenter) Handler(c echo.Context) error {
	if presenterHelper.ShouldEnableInitialSetup(presenter.accountService) {
		return c.Redirect(http.StatusFound, "/setup/")
	}

	rawAccessToken := c.QueryParam("accessToken")
	if rawAccessToken != "" {
		accessToken, err := valueObject.NewAccessTokenStr(rawAccessToken)
		if err == nil {
			sessionCookieExpiresIn := valueObject.NewUnixTimeAfterNow(
				useCase.SessionTokenExpiresIn,
			)
			c.SetCookie(&http.Cookie{
				Name:    infraEnvs.AccessTokenCookieKey,
				Value:   accessToken.String(),
				Path:    "/",
				Expires: sessionCookieExpiresIn.ReadAsGoTime(),
			})
			return c.Redirect(http.StatusFound, "/overview/")
		}
	}

	c.Response().Writer.WriteHeader(http.StatusOK)
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return layoutLogin.Login().
		Render(c.Request().Context(), c.Response().Writer)
}
