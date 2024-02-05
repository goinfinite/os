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
	createAccount dto.CreateAccount,
) error {
	_, err := accQueryRepo.GetByUsername(createAccount.Username)
	if err == nil {
		return errors.New("UsernameAlreadyExists")
	}

	err = accCmdRepo.Add(createAccount)
	if err != nil {
		return errors.New("CreateAccountError")
	}

	log.Printf("Account '%v' created.", createAccount.Username.String())

	return nil
}
