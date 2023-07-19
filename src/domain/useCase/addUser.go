package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/sam/src/domain/dto"
	"github.com/speedianet/sam/src/domain/repository"
)

func AddUser(
	accQueryRepo repository.AccQueryRepo,
	accCmdRepo repository.AccCmdRepo,
	addUser dto.AddUser,
) error {
	_, err := accQueryRepo.GetByUsername(addUser.Username)
	if err == nil {
		return errors.New("UserAlreadyExists")
	}

	err = accCmdRepo.Add(addUser)
	if err != nil {
		return errors.New("AddUserError")
	}

	log.Printf("User '%v' added.", addUser.Username.String())

	return nil
}
