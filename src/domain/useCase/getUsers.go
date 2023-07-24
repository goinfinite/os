package useCase

import (
	"github.com/speedianet/sam/src/domain/entity"
	"github.com/speedianet/sam/src/domain/repository"
)

func GetUsers(
	accQueryRepo repository.AccQueryRepo,
) ([]entity.AccountDetails, error) {
	return accQueryRepo.Get()
}
