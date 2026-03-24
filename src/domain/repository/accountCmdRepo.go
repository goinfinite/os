package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type AccountCmdRepo interface {
	Create(dto.CreateAccount) (tkValueObject.AccountId, error)
	Delete(tkValueObject.AccountId) error
	Update(dto.UpdateAccount) error
	UpdateApiKey(tkValueObject.AccountId) (tkValueObject.AccessTokenValue, error)
	CreateSecureAccessPublicKey(
		dto.CreateSecureAccessPublicKey,
	) (valueObject.SecureAccessPublicKeyId, error)
	DeleteSecureAccessPublicKey(valueObject.SecureAccessPublicKeyId) error
}
