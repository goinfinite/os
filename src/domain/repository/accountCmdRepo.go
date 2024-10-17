package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type AccountCmdRepo interface {
	Create(createAccount dto.CreateAccount) (valueObject.AccountId, error)
	Delete(accountId valueObject.AccountId) error
	UpdatePassword(accountId valueObject.AccountId, password valueObject.Password) error
	UpdateApiKey(accountId valueObject.AccountId) (valueObject.AccessTokenStr, error)
}
