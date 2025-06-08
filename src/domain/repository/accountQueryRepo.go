package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
)

type AccountQueryRepo interface {
	Count(dto.ReadAccountsRequest) (uint64, error)
	Read(dto.ReadAccountsRequest) (dto.ReadAccountsResponse, error)
	ReadFirst(dto.ReadAccountsRequest) (entity.Account, error)
	ReadSecureAccessPublicKeys(
		dto.ReadSecureAccessPublicKeysRequest,
	) (dto.ReadSecureAccessPublicKeysResponse, error)
	ReadFirstSecureAccessPublicKey(
		dto.ReadSecureAccessPublicKeysRequest,
	) (entity.SecureAccessPublicKey, error)
}
