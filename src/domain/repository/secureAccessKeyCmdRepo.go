package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type SecureAccessKeyCmdRepo interface {
	Create(dto.CreateSecureAccessKey) (valueObject.SecureAccessKeyId, error)
	Delete(dto.DeleteSecureAccessKey) error
}
