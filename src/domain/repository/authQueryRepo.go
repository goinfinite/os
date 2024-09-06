package repository

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
)

type AuthQueryRepo interface {
	IsLoginValid(login dto.Login) bool
	ReadAccessTokenDetails(
		token valueObject.AccessTokenStr,
	) (dto.AccessTokenDetails, error)
}
