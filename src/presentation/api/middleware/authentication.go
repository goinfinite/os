package apiMiddleware

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/useCase"
	authInfra "github.com/goinfinite/os/src/infra/auth"
	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	tkInfra "github.com/goinfinite/tk/src/infra"
	"github.com/labstack/echo/v4"
)

func extractAccountIdFromAccessToken(
	authQueryRepo repository.AuthQueryRepo,
	accessTokenStr tkValueObject.AccessTokenValue,
	userIpAddress tkValueObject.IpAddress,
) (accountId tkValueObject.AccountId, err error) {
	trustedCidrs, err := tkInfra.TrustedCidrsReader()
	if err != nil {
		slog.Error("TrustedCidrsReaderError", slog.String("err", err.Error()))
		trustedCidrs = []tkValueObject.CidrBlock{}
	}

	accessTokenDetails, err := useCase.ReadAccessTokenDetails(
		authQueryRepo, accessTokenStr, trustedCidrs, userIpAddress,
	)
	if err != nil {
		return accountId, err
	}

	return accessTokenDetails.AccountId, nil
}

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

			rawAccessToken := ""
			accessTokenCookie, err := echoContext.Cookie(infraEnvs.AccessTokenCookieKey)
			if err == nil {
				rawAccessToken = accessTokenCookie.Value
			}

			if rawAccessToken == "" {
				rawAccessToken = echoContext.Request().Header.Get("Authorization")
				if rawAccessToken == "" {
					return authError("MissingAuthorizationHeader")
				}
				bearerPrefix := "Bearer "
				if !strings.HasPrefix(rawAccessToken, bearerPrefix) {
					return authError("AuthorizationHeaderMissingBearerPrefix")
				}
				rawAccessToken = strings.TrimPrefix(rawAccessToken, bearerPrefix)
			}

			accessTokenStr, err := tkValueObject.NewAccessTokenValue(rawAccessToken)
			if err != nil {
				return authError("InvalidAccessToken")
			}

			userIpAddress, err := tkValueObject.NewIpAddress(echoContext.RealIP())
			if err != nil {
				return authError("InvalidIpAddress")
			}
			authQueryRepo := authInfra.NewAuthQueryRepo(persistentDbSvc)

			accountId, err := extractAccountIdFromAccessToken(
				authQueryRepo, accessTokenStr, userIpAddress,
			)
			if err != nil {
				return authError("InvalidAccessToken")
			}

			echoContext.Set("operatorAccountId", accountId.String())
			return subsequentHandler(echoContext)
		}
	}
}
