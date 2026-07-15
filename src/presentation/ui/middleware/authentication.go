package uiMiddleware

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	authInfra "github.com/goinfinite/os/src/infra/auth"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	sharedHelper "github.com/goinfinite/os/src/presentation/shared/helper"
	"github.com/labstack/echo/v4"
)

func shouldSkipUiAuthentication(
	uiBasePath, apiBasePath string,
	httpRequest *http.Request,
) bool {
	requestPath := httpRequest.URL.Path

	if requestPath == apiBasePath || strings.HasPrefix(requestPath, apiBasePath+"/") {
		return true
	}

	return IsUnauthenticatedUiCall(httpRequest, uiBasePath)
}

func Authentication(
	uiBasePath string,
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) echo.MiddlewareFunc {
	return func(subsequentHandler echo.HandlerFunc) echo.HandlerFunc {
		return func(echoContext echo.Context) error {
			apiBasePath, assertOk := echoContext.Get("apiBasePath").(string)
			if !assertOk {
				slog.Error("AssertApiBasePathFailed")
				return echoContext.NoContent(http.StatusInternalServerError)
			}

			if shouldSkipUiAuthentication(uiBasePath, apiBasePath, echoContext.Request()) {
				return subsequentHandler(echoContext)
			}

			loginPath := uiBasePath + "/login/"

			baseHref, assertOk := echoContext.Get("baseHref").(string)
			if !assertOk {
				return echoContext.NoContent(http.StatusInternalServerError)
			}
			if len(baseHref) > 0 {
				baseHrefNoTrailing := strings.TrimSuffix(baseHref, "/")
				loginPath = baseHrefNoTrailing + loginPath
			}

			authHelper := sharedHelper.AuthenticationHelper{}

			accessToken, err := authHelper.ExtractAccessToken(echoContext)
			if err != nil {
				rawAuthorizationHeader := echoContext.Request().Header.Get(
					sharedHelper.AuthorizationHeader,
				)

				switch {
				case errors.Is(err, sharedHelper.ErrMissingAuthorizationHeader):
					slog.Debug("MissingAuthorizationHeader")
				case errors.Is(err, sharedHelper.ErrAuthorizationHeaderMissingBearerPrefix):
					slog.Debug(
						"AuthorizationHeaderMissingBearerPrefix",
						slog.String("rawAuthorization", rawAuthorizationHeader),
					)
				default:
					slog.Debug(
						"InvalidAccessToken",
						slog.String("rawAccessToken", rawAuthorizationHeader),
					)
				}
				return echoContext.Redirect(http.StatusTemporaryRedirect, loginPath)
			}

			operatorIpAddress, extractionErr := authHelper.ExtractIpAddress(echoContext)
			if extractionErr != nil {
				slog.Debug(
					"InvalidUserIpAddress",
					slog.String("err", extractionErr.Error()),
				)
				return echoContext.Redirect(http.StatusTemporaryRedirect, loginPath)
			}

			authQueryRepo := authInfra.NewAuthQueryRepo(persistentDbSvc)
			_, err = authHelper.ReadAccessTokenAccountId(
				authQueryRepo, accessToken, operatorIpAddress,
			)
			if err != nil {
				slog.Debug("InvalidAccessTokenDetails", slog.String("err", err.Error()))
				return echoContext.Redirect(http.StatusTemporaryRedirect, loginPath)
			}
			return subsequentHandler(echoContext)
		}
	}
}
