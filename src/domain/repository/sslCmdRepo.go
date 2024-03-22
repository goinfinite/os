package repository

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type SslCmdRepo interface {
	Create(createSslPair dto.CreateSslPair) error
	Delete(sslId valueObject.SslId) error
	ReplaceWithValidSsl(sslPair entity.SslPair) error
}
