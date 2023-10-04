package repository

import (
	"github.com/speedianet/sam/src/domain/entity"
)

type SslQueryRepo interface {
	Get() ([]entity.Ssl, error)
}
