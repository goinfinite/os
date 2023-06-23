package repository

import (
	"github.com/speedianet/sam/src/domain/dto"
	"github.com/speedianet/sam/src/domain/valueObject"
)

type AuthQueryRepo interface {
	IsLoginValid(login dto.Login) bool
	GetAccessTokenDetails(
		token valueObject.AccessTokenStr,
	) (dto.AccessTokenDetails, error)
}
