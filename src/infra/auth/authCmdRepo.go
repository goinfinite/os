package authInfra

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
)

type AuthCmdRepo struct {
}

func (repo AuthCmdRepo) CreateSessionToken(
	accountId valueObject.AccountId,
	expiresIn valueObject.UnixTime,
	ipAddress valueObject.IpAddress,
) (entity.AccessToken, error) {
	var accessToken entity.AccessToken

	jwtSecret := os.Getenv("JWT_SECRET")
	apiURL, err := infraHelper.GetPrimaryVirtualHost()
	if err != nil {
		return accessToken, errors.New("PrimaryVirtualHostNotFound")
	}

	now := time.Now()
	tokenExpiration := time.Unix(expiresIn.Int64(), 0)

	claims := jwt.MapClaims{
		"iss":        apiURL,
		"iat":        now.Unix(),
		"nbf":        now.Unix(),
		"exp":        tokenExpiration.Unix(),
		"accountId":  accountId.Uint64(),
		"originalIp": ipAddress.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStrUnparsed, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return accessToken, errors.New("SessionTokenGenerationError")
	}

	tokenType, err := valueObject.NewAccessTokenType("sessionToken")
	if err != nil {
		return accessToken, err
	}

	tokenStr, err := valueObject.NewAccessTokenStr(tokenStrUnparsed)
	if err != nil {
		return accessToken, err
	}

	return entity.NewAccessToken(
		tokenType,
		expiresIn,
		tokenStr,
	), nil
}
