package authInfra

import (
	"errors"
	"os"
	"time"

	"github.com/goinfinite/os/src/domain/entity"
	vhostInfra "github.com/goinfinite/os/src/infra/vhost"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	"github.com/golang-jwt/jwt/v5"
)

type AuthCmdRepo struct {
	sessionTokenSecretBytes []byte
}

func NewAuthCmdRepo() *AuthCmdRepo {
	sessionTokenSecret := os.Getenv("JWT_SECRET")
	return &AuthCmdRepo{
		sessionTokenSecretBytes: []byte(sessionTokenSecret),
	}
}

func (repo *AuthCmdRepo) CreateSessionToken(
	accountId tkValueObject.AccountId,
	expiresIn tkValueObject.UnixTime,
	ipAddress tkValueObject.IpAddress,
) (accessToken entity.AccessToken, err error) {
	apiUrl, err := vhostInfra.NewVirtualHostHelpers().ReadPrimaryVirtualHostHostname()
	if err != nil {
		return accessToken, errors.New("PrimaryVirtualHostNotFound")
	}

	now := time.Now()
	tokenExpiration := time.Unix(expiresIn.Int64(), 0)

	claims := jwt.MapClaims{
		"iss":        apiUrl,
		"iat":        now.Unix(),
		"nbf":        now.Unix(),
		"exp":        tokenExpiration.Unix(),
		"accountId":  accountId.Uint64(),
		"originalIp": ipAddress.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStrUnparsed, err := token.SignedString(repo.sessionTokenSecretBytes)
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

	return entity.NewAccessToken(tokenType, expiresIn, tokenStr), nil
}
