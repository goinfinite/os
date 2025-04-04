package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
)

type SslQueryRepo interface {
	Read(dto.ReadSslPairsRequest) (dto.ReadSslPairsResponse, error)
	ReadFirst(dto.ReadSslPairsRequest) (entity.SslPair, error)
}
