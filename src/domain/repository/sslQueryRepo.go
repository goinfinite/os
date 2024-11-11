package repository

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type SslQueryRepo interface {
	Read() ([]entity.SslPair, error)
	ReadById(sslPairId valueObject.SslPairId) (entity.SslPair, error)
}
