package authInfra

import (
	"errors"
	"os"
	"time"

	"github.com/goinfinite/os/src/domain/entity"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	"github.com/golang-jwt/jwt"
)

type AuthCmdRepo struct {
}

func (repo AuthCmdRepo) CreateSessionToken(
	accountId tkValueObject.AccountId,
	expiresIn tkValueObject.UnixTime,
	ipAddress tkValueObject.IpAddress,
) (entity.AccessToken, error) {
	var accessToken entity.AccessToken

	jwtSecret := os.Getenv("JWT_SECRET")
	apiURL, err := infraHelper.ReadPrimaryVirtualHostHostname()
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

	tokenType, err := tkValueObject.NewAccessTokenType("sessionToken")
	if err != nil {
		return accessToken, err
	}

	tokenStr, err := tkValueObject.NewAccessTokenValue(tokenStrUnparsed)
	if err != nil {
		return accessToken, err
	}

	return entity.NewAccessToken(
		tokenType,
		expiresIn,
		tokenStr,
	), nil
}
