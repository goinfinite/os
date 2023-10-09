package useCase

import (
	"github.com/speedianet/sam/src/domain/entity"
	"github.com/speedianet/sam/src/domain/repository"
)

func GetSsls(
	sslQueryRepo repository.SslQueryRepo,
) ([]entity.Ssl, error) {
	return sslQueryRepo.Get()
}
