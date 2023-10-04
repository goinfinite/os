package repository

import (
	"github.com/speedianet/sam/src/domain/dto"
)

type SslCmdRepo interface {
	Add(addSsl dto.AddSsl) error
}
