package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type SslCmdRepo interface {
	Create(createSslPair dto.CreateSslPair) error
	Delete(sslId valueObject.SslId) error
	ReplaceWithValidSsl(sslPair entity.SslPair) error
	DeleteSslPairVhosts(deleteDto dto.DeleteSslPairVhosts) error
}
