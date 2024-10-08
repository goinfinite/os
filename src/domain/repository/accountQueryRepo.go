package repository

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type AccountQueryRepo interface {
	Read() ([]entity.Account, error)
	ReadByUsername(username valueObject.Username) (entity.Account, error)
	ReadById(accountId valueObject.AccountId) (entity.Account, error)
}
