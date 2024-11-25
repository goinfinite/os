package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type AccountCmdRepo interface {
	Create(dto.CreateAccount) (valueObject.AccountId, error)
	Delete(valueObject.AccountId) error
	UpdatePassword(valueObject.AccountId, valueObject.Password) error
	UpdateApiKey(valueObject.AccountId) (valueObject.AccessTokenStr, error)
	CreateSecureAccessKey(
		dto.CreateSecureAccessKey,
	) (valueObject.SecureAccessKeyId, error)
	DeleteSecureAccessKey(dto.DeleteSecureAccessKey) error
}
