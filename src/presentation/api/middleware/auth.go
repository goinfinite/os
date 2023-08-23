package apiMiddleware

import (
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/sam/src/domain/useCase"
	"github.com/speedianet/sam/src/domain/valueObject"
	"github.com/speedianet/sam/src/infra"
)

func getAccountIdFromAccessToken(
	accessToken valueObject.AccessTokenStr,
	ipAddress valueObject.IpAddress,
) (valueObject.AccountId, error) {
	authQueryRepo := infra.AuthQueryRepo{}

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
		accessToken,
		trustedIps,
		ipAddress,
	)
	if err != nil {
		return valueObject.AccountId(0), err
	}

	return accessTokenDetails.AccountId, nil
}

func Auth(basePath string) echo.MiddlewareFunc {
	urlSkipRegex := regexp.MustCompile(
		"^" + basePath + "/" + "(swagger|auth|health)",
	)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if urlSkipRegex.MatchString(c.Request().URL.Path) {
				return next(c)
			}

			token := c.Request().Header.Get("Authorization")
			if token == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, map[string]interface{}{
					"status": http.StatusUnauthorized,
					"body":   "MissingAuthToken",
				})
			}

			tokenWithoutPrefix := token[7:]
			accountId, err := getAccountIdFromAccessToken(
				valueObject.AccessTokenStr(tokenWithoutPrefix),
				valueObject.NewIpAddressPanic(c.RealIP()),
			)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, map[string]interface{}{
					"status": http.StatusUnauthorized,
					"body":   "InvalidAuthToken",
				})
			}

			c.Set("accountId", accountId.String())
			return next(c)
		}
	}
}
