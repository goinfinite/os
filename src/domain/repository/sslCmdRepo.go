package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type SslCmdRepo interface {
	Create(dto.CreateSslPair) (valueObject.SslPairId, error)
	Delete(valueObject.SslPairId) error
	CreatePubliclyTrusted(dto.CreatePubliclyTrustedSslPair) error
	DeleteSslPairVhosts(dto.DeleteSslPairVhosts) error
}
