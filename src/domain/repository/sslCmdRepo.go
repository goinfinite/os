package repository

import (
	"math/big"

	"github.com/speedianet/sam/src/domain/dto"
)

type SslCmdRepo interface {
	Add(addSsl dto.AddSsl) error
	Delete(sslId *big.Int) error
}
