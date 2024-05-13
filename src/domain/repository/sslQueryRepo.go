package repository

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type SslQueryRepo interface {
	Read() ([]entity.SslPair, error)
	ReadById(sslId valueObject.SslId) (entity.SslPair, error)
}
