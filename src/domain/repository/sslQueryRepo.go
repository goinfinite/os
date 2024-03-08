package repository

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type SslQueryRepo interface {
	GetSslPairs() ([]entity.SslPair, error)
	GetSslPairById(sslId valueObject.SslId) (entity.SslPair, error)
	GetSslPairByHostname(vhost valueObject.Fqdn) (entity.SslPair, error)
	GetOwnershipHash(sslCrtContent valueObject.SslCertificateContent) string
	IsSslPairValid(vhost valueObject.Fqdn) bool
	ValidateSslOwnership(vhost valueObject.Fqdn, ownershipHash string) bool
}
