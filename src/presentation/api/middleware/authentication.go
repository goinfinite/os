package apiMiddleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	authInfra "github.com/goinfinite/os/src/infra/auth"
	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	"github.com/labstack/echo/v4"
)

func extractAccountIdFromAccessToken(
	authQueryRepo repository.AuthQueryRepo,
	accessTokenStr valueObject.AccessTokenStr,
	ipAddress valueObject.IpAddress,
) (valueObject.AccountId, error) {
	trustedIpsRaw := strings.Split(os.Getenv("TRUSTED_IPS"), ",")
	var trustedIps []valueObject.IpAddress
	for _, trustedIp := range trustedIpsRaw {
		ipAddress, err := valueObject.NewIpAddress(trustedIp)
		if err != nil {
			continue
		}
		trustedIps = append(trustedIps, ipAddress)
	}

	accessTokenDetails, err := useCase.ReadAccessTokenDetails(
		authQueryRepo, accessTokenStr, trustedIps, ipAddress,
	)
	if err != nil {
		return valueObject.AccountId(0), err
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
					return authError("MissingAccessToken")
				}
				tokenWithoutPrefix := rawAccessToken[7:]
				rawAccessToken = tokenWithoutPrefix
			}

			accessTokenStr, err := valueObject.NewAccessTokenStr(rawAccessToken)
			if err != nil {
				return authError("InvalidAccessToken")
			}

			userIpAddress, err := valueObject.NewIpAddress(echoContext.RealIP())
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
