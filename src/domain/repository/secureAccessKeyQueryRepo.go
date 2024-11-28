package repository

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type SecureAccessKeyQueryRepo interface {
	Read(valueObject.AccountId) ([]entity.SecureAccessKey, error)
}
