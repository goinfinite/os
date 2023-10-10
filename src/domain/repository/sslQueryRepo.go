package repository

import (
	"github.com/speedianet/sam/src/domain/entity"
	"github.com/speedianet/sam/src/domain/valueObject"
)

type SslQueryRepo interface {
	Get() ([]entity.SslPair, error)
	GetById(sslSerialNumber valueObject.SslSerialNumber) (entity.SslPair, error)
}
