package uiPresenter

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	"github.com/goinfinite/os/src/presentation/liaison"
	layoutLogin "github.com/goinfinite/os/src/presentation/ui/layout/login"
	presenterHelper "github.com/goinfinite/os/src/presentation/ui/presenter/helper"
	"github.com/labstack/echo/v4"
)

type LoginPresenter struct {
	accountLiaison *liaison.AccountLiaison
}

func NewLoginPresenter(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *LoginPresenter {
	return &LoginPresenter{
		accountLiaison: liaison.NewAccountLiaison(persistentDbSvc, trailDbSvc),
	}
}

func (presenter *LoginPresenter) Handler(echoContext echo.Context) error {
	uiBasePath, assertOk := echoContext.Get("uiBasePath").(string)
	if !assertOk {
		slog.Error("AssertUiBasePathFailed")
		return echoContext.NoContent(http.StatusInternalServerError)
	}

	baseHref, assertOk := echoContext.Get("baseHref").(string)
	if !assertOk {
		slog.Error("AssertBaseHrefFailed")
		return echoContext.NoContent(http.StatusInternalServerError)
	}
	if len(baseHref) > 0 {
		baseHrefNoTrailing := strings.TrimSuffix(baseHref, "/")
		uiBasePath = baseHrefNoTrailing + uiBasePath
	}

	if presenterHelper.ShouldEnableInitialSetup(presenter.accountLiaison) {
		return echoContext.Redirect(http.StatusFound, uiBasePath+"/setup/")
	}

	rawAccessToken := echoContext.QueryParam("accessToken")
	if rawAccessToken != "" {
		accessToken, err := valueObject.NewAccessTokenStr(rawAccessToken)
		if err == nil {
			sessionCookieExpiresIn := valueObject.NewUnixTimeAfterNow(
				useCase.SessionTokenExpiresIn,
			)
			cookiePathStr := uiBasePath
			if len(cookiePathStr) == 0 {
				cookiePathStr = "/"
			}

			echoContext.SetCookie(&http.Cookie{
				Name:     infraEnvs.AccessTokenCookieKey,
				Value:    accessToken.String(),
				Path:     cookiePathStr,
				Expires:  sessionCookieExpiresIn.ReadAsGoTime(),
				MaxAge:   int(useCase.SessionTokenExpiresIn.Seconds()),
				HttpOnly: false,
				Secure:   true,
				SameSite: http.SameSiteLaxMode,
			})
			return echoContext.Redirect(http.StatusFound, uiBasePath+"/overview/")
		}
		slog.Debug("InvalidAccessTokenDetails", slog.String("rawAccessToken", rawAccessToken))
	}

	loginLayoutSettings := layoutLogin.LoginLayoutSettings{BaseHref: baseHref}

	rawPrefilledUsername := echoContext.QueryParam("prefilledUsername")
	if rawPrefilledUsername != "" {
		username, err := valueObject.NewUsername(rawPrefilledUsername)
		if err == nil {
			loginLayoutSettings.PrefilledUsername = username.String()
		}
	}

	rawPrefilledPassword := echoContext.QueryParam("prefilledPassword")
	if rawPrefilledPassword != "" {
		password, err := valueObject.NewPassword(rawPrefilledPassword)
		if err == nil {
			loginLayoutSettings.PrefilledPassword = password.String()
		}
	}

	echoContext.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	echoContext.Response().Writer.WriteHeader(http.StatusOK)

	return layoutLogin.Login(loginLayoutSettings).
		Render(echoContext.Request().Context(), echoContext.Response().Writer)
}
