package repository

import (
	"github.com/goinfinite/os/src/domain/entity"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type AuthCmdRepo interface {
	CreateSessionToken(
		accountId tkValueObject.AccountId,
		expiresIn tkValueObject.UnixTime,
		ipAddress tkValueObject.IpAddress,
	) (entity.AccessToken, error)
}
