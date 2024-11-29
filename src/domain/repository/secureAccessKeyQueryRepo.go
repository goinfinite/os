package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
)

type SecureAccessKeyQueryRepo interface {
	Read(dto.ReadSecureAccessKeysRequest) (dto.ReadSecureAccessKeysResponse, error)
	ReadFirst(dto.ReadSecureAccessKeysRequest) (entity.SecureAccessKey, error)
}
