package repository

import (
	"github.com/speedianet/sam/src/domain/dto"
	"github.com/speedianet/sam/src/domain/valueObject"
)

type SslCmdRepo interface {
	Add(addSslPair dto.AddSslPair) error
	Delete(sslId valueObject.SslId) error
}
