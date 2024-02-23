package repository

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type AuthCmdRepo interface {
	GenerateSessionToken(
		accountId valueObject.AccountId,
		expiresIn valueObject.UnixTime,
		ipAddress valueObject.IpAddress,
	) (entity.AccessToken, error)
}
