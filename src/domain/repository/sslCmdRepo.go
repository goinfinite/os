package repository

import (
	"github.com/speedianet/sam/src/domain/dto"
	"github.com/speedianet/sam/src/domain/valueObject"
)

type SslCmdRepo interface {
	Add(addSsl dto.AddSsl) error
	Delete(sslId valueObject.SslId) error
}
