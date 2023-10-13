package useCase

import (
	"github.com/speedianet/sam/src/domain/entity"
	"github.com/speedianet/sam/src/domain/repository"
)

func GetSslPairs(
	sslQueryRepo repository.SslQueryRepo,
) ([]entity.SslPair, error) {
	return sslQueryRepo.GetSslPairs()
}
