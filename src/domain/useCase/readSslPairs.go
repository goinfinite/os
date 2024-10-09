package useCase

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/repository"
)

func ReadSslPairs(
	sslQueryRepo repository.SslQueryRepo,
) ([]entity.SslPair, error) {
	return sslQueryRepo.Read()
}
