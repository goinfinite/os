package uiMiddleware

import (
	"log/slog"
	"net/http"
	"os"
	"regexp"
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

func shouldSkipUiAuthentication(
	uiBasePath, apiBasePath string,
	httpRequest *http.Request,
) bool {
	urlSkipRegex := regexp.MustCompile(
		"^(" + apiBasePath + "|" + uiBasePath + "/(login|assets|setup))/",
	)
	return urlSkipRegex.MatchString(httpRequest.URL.Path)
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

			rawAccessToken := ""
			accessTokenCookie, err := echoContext.Cookie(infraEnvs.AccessTokenCookieKey)
			if err == nil {
				rawAccessToken = accessTokenCookie.Value
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

			if rawAccessToken == "" {
				rawAccessToken = echoContext.Request().Header.Get("Authorization")
				if rawAccessToken == "" {
					slog.Debug("MissingAuthorizationHeader")
					return echoContext.Redirect(http.StatusTemporaryRedirect, loginPath)
				}
				bearerPrefix := "Bearer "
				if !strings.HasPrefix(rawAccessToken, bearerPrefix) {
					slog.Debug(
						"AuthorizationHeaderMissingBearerPrefix",
						slog.String("rawAuthorization", rawAccessToken),
					)
					return echoContext.Redirect(http.StatusTemporaryRedirect, loginPath)
				}
				rawAccessToken = strings.TrimPrefix(rawAccessToken, bearerPrefix)
			}

			accessTokenStr, err := valueObject.NewAccessTokenStr(rawAccessToken)
			if err != nil {
				slog.Debug("InvalidAccessToken", slog.String("rawAccessToken", rawAccessToken))
				return echoContext.Redirect(http.StatusTemporaryRedirect, loginPath)
			}

			userIpAddress, err := valueObject.NewIpAddress(echoContext.RealIP())
			if err != nil {
				slog.Debug("InvalidUserIpAddress", slog.String("ipAddress", echoContext.RealIP()))
				return echoContext.Redirect(http.StatusTemporaryRedirect, loginPath)
			}

			authQueryRepo := authInfra.NewAuthQueryRepo(persistentDbSvc)
			_, err = extractAccountIdFromAccessToken(
				authQueryRepo, accessTokenStr, userIpAddress,
			)
			if err != nil {
				slog.Debug("InvalidAccessTokenDetails", slog.String("err", err.Error()))
				return echoContext.Redirect(http.StatusTemporaryRedirect, loginPath)
			}
			return subsequentHandler(echoContext)
		}
	}
}
