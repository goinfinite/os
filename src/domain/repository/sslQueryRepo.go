package repository

import (
	"github.com/speedianet/sam/src/domain/entity"
	"github.com/speedianet/sam/src/domain/valueObject"
)

type SslQueryRepo interface {
	Get() ([]entity.Ssl, error)
	GetById(sslId valueObject.SslId) (entity.Ssl, error)
}
