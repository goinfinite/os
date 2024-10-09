package useCase

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/repository"
)

func GetAccounts(
	accQueryRepo repository.AccQueryRepo,
) ([]entity.Account, error) {
	return accQueryRepo.Get()
}
