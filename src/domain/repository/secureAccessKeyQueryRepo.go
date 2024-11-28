package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
)

type SecureAccessKeyQueryRepo interface {
	Read(dto.ReadSecureAccessKeysRequest) (dto.ReadSecureAccessKeysResponse, error)
}
