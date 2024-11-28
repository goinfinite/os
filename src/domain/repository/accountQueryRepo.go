package repository

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type AccountQueryRepo interface {
	Read() ([]entity.Account, error)
	ReadByUsername(valueObject.Username) (entity.Account, error)
	ReadById(valueObject.AccountId) (entity.Account, error)
}
