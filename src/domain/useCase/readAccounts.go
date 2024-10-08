package useCase

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/repository"
)

func ReadAccounts(
	accountQueryRepo repository.AccountQueryRepo,
) ([]entity.Account, error) {
	return accountQueryRepo.Read()
}
