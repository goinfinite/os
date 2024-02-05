package authInfra

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
)

type AuthCmdRepo struct {
}

func (repo AuthCmdRepo) GenerateSessionToken(
	accountId valueObject.AccountId,
	expiresIn valueObject.UnixTime,
	ipAddress valueObject.IpAddress,
) entity.AccessToken {
	jwtSecret := os.Getenv("JWT_SECRET")
	apiURL, err := infraHelper.GetPrimaryHostname()
	if err != nil {
		panic("PrimaryHostnameNotFound")
	}

	now := time.Now()
	tokenExpiration := time.Unix(expiresIn.Get(), 0)

	claims := jwt.MapClaims{
		"iss":        apiURL,
		"iat":        now.Unix(),
		"nbf":        now.Unix(),
		"exp":        tokenExpiration.Unix(),
		"accountId":  accountId.Get(),
		"originalIp": ipAddress.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStrUnparsed, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		panic("SessionTokenGenerationError")
	}

	tokenType := valueObject.NewAccessTokenTypePanic("sessionToken")
	tokenStr := valueObject.NewAccessTokenStrPanic(tokenStrUnparsed)

	return entity.NewAccessToken(
		tokenType,
		expiresIn,
		tokenStr,
	)
}
