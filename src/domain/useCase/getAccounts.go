package useCase

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/repository"
)

func GetAccounts(
	accQueryRepo repository.AccQueryRepo,
) ([]entity.Account, error) {
	return accQueryRepo.Get()
}
