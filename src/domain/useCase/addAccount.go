package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/sam/src/domain/dto"
	"github.com/speedianet/sam/src/domain/repository"
)

func AddAccount(
	accQueryRepo repository.AccQueryRepo,
	accCmdRepo repository.AccCmdRepo,
	addAccount dto.AddAccount,
) error {
	_, err := accQueryRepo.GetByUsername(addAccount.Username)
	if err == nil {
		return errors.New("UsernameAlreadyExists")
	}

	err = accCmdRepo.Add(addAccount)
	if err != nil {
		return errors.New("AddAccountError")
	}

	log.Printf("Account '%v' added.", addAccount.Username.String())

	return nil
}
