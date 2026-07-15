package sharedHelper

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/useCase"
	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	tkInfra "github.com/goinfinite/tk/src/infra"
	tkPresentation "github.com/goinfinite/tk/src/presentation"
	"github.com/labstack/echo/v4"
)

const AuthorizationHeader string = "Authorization"

var (
	ErrMissingAuthorizationHeader error = errors.New(
		"MissingAuthorizationHeader",
	)
	ErrAuthorizationHeaderMissingBearerPrefix error = errors.New(
		"AuthorizationHeaderMissingBearerPrefix",
	)
)

type AuthenticationHelper struct{}

func (AuthenticationHelper) ExtractAccessToken(
	echoContext echo.Context,
) (accessToken tkValueObject.AccessTokenValue, err error) {
	rawAccessToken := ""
	accessTokenCookie, cookieErr := echoContext.Cookie(infraEnvs.AccessTokenCookieKey)
	if cookieErr == nil {
		rawAccessToken = accessTokenCookie.Value
	}

	if rawAccessToken == "" {
		rawAccessToken = echoContext.Request().Header.Get(AuthorizationHeader)
		if rawAccessToken == "" {
			return accessToken, ErrMissingAuthorizationHeader
		}
		bearerPrefix := "Bearer "
		if !strings.HasPrefix(rawAccessToken, bearerPrefix) {
			return accessToken, ErrAuthorizationHeaderMissingBearerPrefix
		}
		rawAccessToken = strings.TrimPrefix(rawAccessToken, bearerPrefix)
	}

	return tkValueObject.NewAccessTokenValue(rawAccessToken)
}

func (AuthenticationHelper) ExtractIpAddress(
	echoContext echo.Context,
) (tkValueObject.IpAddress, error) {
	return tkPresentation.NewRequesterIpExtractor().Execute(
		echoContext.Request(),
	)
}

func (AuthenticationHelper) ReadAccessTokenAccountId(
	authQueryRepo repository.AuthQueryRepo,
	accessToken tkValueObject.AccessTokenValue,
	userIpAddress tkValueObject.IpAddress,
) (accountId tkValueObject.AccountId, err error) {
	trustedCidrs, trustedCidrsErr := tkInfra.TrustedCidrsReader()
	if trustedCidrsErr != nil {
		slog.Error(
			"TrustedCidrsReaderError",
			slog.String("err", trustedCidrsErr.Error()),
		)
		trustedCidrs = []tkValueObject.CidrBlock{}
	}

	accessTokenDetails, err := useCase.ReadAccessTokenDetails(
		authQueryRepo, accessToken, trustedCidrs, userIpAddress,
	)
	if err != nil {
		return accountId, err
	}

	return accessTokenDetails.AccountId, nil
}
