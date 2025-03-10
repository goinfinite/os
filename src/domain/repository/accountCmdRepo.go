package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type AccountCmdRepo interface {
	Create(dto.CreateAccount) (valueObject.AccountId, error)
	Delete(valueObject.AccountId) error
	Update(dto.UpdateAccount) error
	UpdateApiKey(valueObject.AccountId) (valueObject.AccessTokenStr, error)
	CreateSecureAccessPublicKey(
		dto.CreateSecureAccessPublicKey,
	) (valueObject.SecureAccessPublicKeyId, error)
	DeleteSecureAccessPublicKey(valueObject.SecureAccessPublicKeyId) error
}
