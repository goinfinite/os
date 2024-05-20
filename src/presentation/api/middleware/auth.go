package apiMiddleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	authInfra "github.com/speedianet/os/src/infra/auth"
)

func getAccountIdFromAccessToken(
	authQueryRepo repository.AuthQueryRepo,
	accessTokenValue valueObject.AccessTokenStr,
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

	accessTokenDetails, err := useCase.GetAccessTokenDetails(
		authQueryRepo,
		accessTokenValue,
		trustedIps,
		ipAddress,
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

func Auth(apiBasePath string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			shouldSkip := IsSkippableApiCall(c.Request(), apiBasePath)
			if shouldSkip {
				return next(c)
			}

			rawAccessToken := ""
			accessTokenCookie, err := c.Cookie("os-auth-token")
			if err == nil {
				rawAccessToken = accessTokenCookie.Value
			}

			if rawAccessToken == "" {
				rawAccessToken = c.Request().Header.Get("Authorization")
				if rawAccessToken == "" {
					return authError("MissingAccessToken")
				}
				tokenWithoutPrefix := rawAccessToken[7:]
				rawAccessToken = tokenWithoutPrefix
			}

			accessTokenValue, err := valueObject.NewAccessTokenStr(rawAccessToken)
			if err != nil {
				return authError("InvalidAccessToken")
			}

			userIpAddress, err := valueObject.NewIpAddress(c.RealIP())
			if err != nil {
				return authError("InvalidIpAddress")
			}

			authQueryRepo := authInfra.AuthQueryRepo{}
			accountId, err := getAccountIdFromAccessToken(
				authQueryRepo, accessTokenValue, userIpAddress,
			)
			if err != nil {
				return authError("InvalidAccessToken")
			}

			c.Set("accountId", accountId.String())
			return next(c)
		}
	}
}
