package presenter

import (
	"net/http"

	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	"github.com/goinfinite/os/src/presentation/ui/layout"
	"github.com/labstack/echo/v4"
)

type LoginPresenter struct {
}

func NewLoginPresenter() *LoginPresenter {
	return &LoginPresenter{}
}

func (presenter *LoginPresenter) Handler(c echo.Context) error {
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
			return c.Redirect(http.StatusFound, "/_/#/overview")
		}
	}

	c.Response().Writer.WriteHeader(http.StatusOK)
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return layout.Login().
		Render(c.Request().Context(), c.Response().Writer)
}
