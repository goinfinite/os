package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type AuthQueryRepo interface {
	IsLoginValid(createDto dto.CreateSessionToken) bool
	ReadAccessTokenDetails(
		token tkValueObject.AccessTokenValue,
	) (dto.AccessTokenDetails, error)
}
