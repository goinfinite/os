package useCase

import (
	"github.com/speedianet/sam/src/domain/entity"
	"github.com/speedianet/sam/src/domain/repository"
)

func GetAccounts(
	accQueryRepo repository.AccQueryRepo,
) ([]entity.Account, error) {
	return accQueryRepo.Get()
}
