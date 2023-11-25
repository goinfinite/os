package repository

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type AccQueryRepo interface {
	Get() ([]entity.Account, error)
	GetByUsername(
		username valueObject.Username,
	) (entity.Account, error)
	GetById(
		accId valueObject.AccountId,
	) (entity.Account, error)
}
