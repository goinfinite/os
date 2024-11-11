package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type SslCmdRepo interface {
	Create(createSslPair dto.CreateSslPair) (valueObject.SslPairId, error)
	Delete(sslPairId valueObject.SslPairId) error
	ReplaceWithValidSsl(sslPair entity.SslPair) error
	DeleteSslPairVhosts(deleteDto dto.DeleteSslPairVhosts) error
}
