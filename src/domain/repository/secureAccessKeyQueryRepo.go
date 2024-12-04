package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
)

type SecureAccessKeyQueryRepo interface {
	Read(dto.ReadSecureAccessPublicKeysRequest) (dto.ReadSecureAccessPublicKeysResponse, error)
	ReadFirst(dto.ReadSecureAccessPublicKeysRequest) (entity.SecureAccessPublicKey, error)
}
