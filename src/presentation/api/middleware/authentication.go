package apiMiddleware

import (
	"errors"
	"log/slog"
	"net/http"

	authInfra "github.com/goinfinite/os/src/infra/auth"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	sharedHelper "github.com/goinfinite/os/src/presentation/shared/helper"
	"github.com/labstack/echo/v4"
)

func authError(message string) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusUnauthorized, map[string]interface{}{
		"status": http.StatusUnauthorized,
		"body":   message,
	})
}

func Authentication(
	apiBasePath string,
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) echo.MiddlewareFunc {
	return func(subsequentHandler echo.HandlerFunc) echo.HandlerFunc {
		return func(echoContext echo.Context) error {
			shouldSkip := IsSkippableApiCall(echoContext.Request(), apiBasePath)
			if shouldSkip {
				return subsequentHandler(echoContext)
			}

			authHelper := sharedHelper.AuthenticationHelper{}

			accessToken, err := authHelper.ExtractAccessToken(echoContext)
			if err != nil {
				switch {
				case errors.Is(err, sharedHelper.ErrMissingAuthorizationHeader):
					return authError("MissingAuthorizationHeader")
				case errors.Is(err, sharedHelper.ErrAuthorizationHeaderMissingBearerPrefix):
					return authError("AuthorizationHeaderMissingBearerPrefix")
				default:
					return authError("InvalidAccessToken")
				}
			}

			operatorIpAddress, extractionErr := authHelper.ExtractIpAddress(echoContext)
			if extractionErr != nil {
				slog.Debug(
					"InvalidOperatorIpAddress",
					slog.String("source", "AuthenticationMiddleware"),
					slog.String("err", extractionErr.Error()),
				)
				return authError("InvalidOperatorIpAddress")
			}

			authQueryRepo := authInfra.NewAuthQueryRepo(persistentDbSvc)

			accountId, err := authHelper.ReadAccessTokenAccountId(
				authQueryRepo, accessToken, operatorIpAddress,
			)
			if err != nil {
				slog.Debug(
					"InvalidAccessToken",
					slog.String("source", "AuthenticationMiddleware"),
					slog.String("err", err.Error()),
				)
				return authError("InvalidAccessToken")
			}

			echoContext.Set("operatorAccountId", accountId)
			echoContext.Set("operatorIpAddress", operatorIpAddress)
			return subsequentHandler(echoContext)
		}
	}
}
