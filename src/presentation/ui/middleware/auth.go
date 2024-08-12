package uiMiddleware

import (
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	authInfra "github.com/speedianet/os/src/infra/auth"
)

func getAccountIdFromAccessToken(
	authQueryRepo repository.AuthQueryRepo,
	accessToken valueObject.AccessTokenStr,
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
		authQueryRepo, accessToken, trustedIps, ipAddress,
	)
	if err != nil {
		return valueObject.AccountId(0), err
	}

	return accessTokenDetails.AccountId, nil
}

func shouldSkipUiAuthentication(req *http.Request) bool {
	urlSkipRegex := regexp.MustCompile(`^/(api|\_|login)/`)
	return urlSkipRegex.MatchString(req.URL.Path)
}

func Authentication() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if shouldSkipUiAuthentication(c.Request()) {
				return next(c)
			}

			rawAccessToken := ""
			accessTokenCookie, err := c.Cookie("control-access-token")
			if err == nil {
				rawAccessToken = accessTokenCookie.Value
			}

			loginPath := "/_/#/login"

			if rawAccessToken == "" {
				rawAccessToken = c.Request().Header.Get("Authorization")
				if rawAccessToken == "" {
					return c.Redirect(http.StatusTemporaryRedirect, loginPath)
				}
				tokenWithoutPrefix := rawAccessToken[7:]
				rawAccessToken = tokenWithoutPrefix
			}

			accessToken, err := valueObject.NewAccessTokenStr(rawAccessToken)
			if err != nil {
				return c.Redirect(http.StatusTemporaryRedirect, loginPath)
			}

			userIpAddress, err := valueObject.NewIpAddress(c.RealIP())
			if err != nil {
				return c.Redirect(http.StatusTemporaryRedirect, loginPath)
			}

			authQueryRepo := authInfra.AuthQueryRepo{}
			_, err = getAccountIdFromAccessToken(
				authQueryRepo, accessToken, userIpAddress,
			)
			if err != nil {
				return c.Redirect(http.StatusTemporaryRedirect, loginPath)
			}
			return next(c)
		}
	}
}
