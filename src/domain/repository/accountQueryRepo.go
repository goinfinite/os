package repository

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type AccountQueryRepo interface {
	Read() ([]entity.Account, error)
	ReadByUsername(username valueObject.Username) (entity.Account, error)
	ReadById(accountId valueObject.AccountId) (entity.Account, error)
	ReadSecureAccessKeys(
		accountId valueObject.AccountId,
	) ([]entity.SecureAccessKey, error)
}
