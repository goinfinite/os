package repository

import (
	"math/big"

	"github.com/speedianet/sam/src/domain/entity"
)

type SslQueryRepo interface {
	Get() ([]entity.Ssl, error)
	GetById(sslId *big.Int) (entity.Ssl, error)
}
