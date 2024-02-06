package repository

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
)

type SslCmdRepo interface {
	GenerateSelfSignedCert(vhost valueObject.Fqdn) error
	Add(addSslPair dto.AddSslPair) error
	Delete(sslId valueObject.SslId) error
}
