package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
)

func CreateAccount(
	accQueryRepo repository.AccQueryRepo,
	accCmdRepo repository.AccCmdRepo,
	addAccount dto.CreateAccount,
) error {
	_, err := accQueryRepo.GetByUsername(addAccount.Username)
	if err == nil {
		return errors.New("UsernameAlreadyExists")
	}

	err = accCmdRepo.Add(addAccount)
	if err != nil {
		return errors.New("CreateAccountError")
	}

	log.Printf("Account '%v' added.", addAccount.Username.String())

	return nil
}
