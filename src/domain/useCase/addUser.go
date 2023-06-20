package useCase

import (
	"log"

	"github.com/speedianet/sam/src/domain/dto"
	"github.com/speedianet/sam/src/domain/repository"
)

func AddUser(
	accQueryRepo repository.AccQueryRepo,
	accCmdRepo repository.AccCmdRepo,
	addUser dto.AddUser,
) {
	_, err := accQueryRepo.GetByUsername(addUser.Username)
	if err == nil {
		panic("UsernameAlreadyExists")
	}

	err = accCmdRepo.Add(addUser)
	if err != nil {
		panic("UserAddError")
	}

	log.Printf("User '%v' added.", addUser.Username.String())
}
