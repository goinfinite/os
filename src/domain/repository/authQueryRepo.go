package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type AuthQueryRepo interface {
	IsLoginValid(createDto dto.CreateSessionToken) bool
	ReadAccessTokenDetails(
		token valueObject.AccessTokenStr,
	) (dto.AccessTokenDetails, error)
}
