package repository

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type SslQueryRepo interface {
	GetSslPairs() ([]entity.SslPair, error)
	GetSslPairById(sslId valueObject.SslId) (entity.SslPair, error)
	GetSslPairByVirtualHost(virtualHost valueObject.Fqdn) (entity.SslPair, error)
}
