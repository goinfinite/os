package repository

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
)

type SslCmdRepo interface {
	Add(addSslPair dto.CreateSslPair) error
	Delete(sslId valueObject.SslId) error
}
