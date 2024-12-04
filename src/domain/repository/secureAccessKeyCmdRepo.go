package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type SecureAccessKeyCmdRepo interface {
	Create(dto.CreateSecureAccessPublicKey) (valueObject.SecureAccessPublicKeyId, error)
	Delete(valueObject.SecureAccessPublicKeyId) error
}
