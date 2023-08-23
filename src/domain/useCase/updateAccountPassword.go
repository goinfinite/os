package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/sam/src/domain/dto"
	"github.com/speedianet/sam/src/domain/repository"
)

func UpdateAccountPassword(
	accQueryRepo repository.AccQueryRepo,
	accCmdRepo repository.AccCmdRepo,
	updateAccountDto dto.UpdateAccount,
) error {
	_, err := accQueryRepo.GetById(updateAccountDto.AccountId)
	if err != nil {
		return errors.New("AccountNotFound")
	}

	err = accCmdRepo.UpdatePassword(
		updateAccountDto.AccountId,
		*updateAccountDto.Password,
	)
	if err != nil {
		return errors.New("UpdateAccountPasswordError")
	}

	log.Printf("AccountId '%v' password updated.", updateAccountDto.AccountId)

	return nil
}
