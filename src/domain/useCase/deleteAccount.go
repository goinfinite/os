package useCase

import (
	"errors"
	"log"

	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
)

func DeleteAccount(
	accQueryRepo repository.AccQueryRepo,
	accCmdRepo repository.AccCmdRepo,
	accountId valueObject.AccountId,
) error {
	_, err := accQueryRepo.GetById(accountId)
	if err != nil {
		return errors.New("AccountNotFound")
	}

	err = accCmdRepo.Delete(accountId)
	if err != nil {
		return errors.New("DeleteAccountError")
	}

	log.Printf("AccountId '%v' deleted.", accountId)

	return nil
}
