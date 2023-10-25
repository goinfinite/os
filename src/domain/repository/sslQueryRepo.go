package repository

import (
	"github.com/speedianet/sam/src/domain/entity"
	"github.com/speedianet/sam/src/domain/valueObject"
)

type SslQueryRepo interface {
	GetSslPairs() ([]entity.SslPair, error)
	GetSslPairById(sslId valueObject.SslId) (entity.SslPair, error)
}
