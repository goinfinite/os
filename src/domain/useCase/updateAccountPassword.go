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
	_, err := accQueryRepo.GetById(updateAccountDto.Id)
	if err != nil {
		return errors.New("AccountNotFound")
	}

	err = accCmdRepo.UpdatePassword(
		updateAccountDto.Id,
		*updateAccountDto.Password,
	)
	if err != nil {
		return errors.New("UpdateAccountPasswordError")
	}

	log.Printf("AccountId '%v' password updated.", updateAccountDto.Id)

	return nil
}
