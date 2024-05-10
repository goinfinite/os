package useCase

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/repository"
)

func ReadSslPairs(
	sslQueryRepo repository.SslQueryRepo,
) ([]entity.SslPair, error) {
	return sslQueryRepo.Read()
}
