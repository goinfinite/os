package useCase

import (
	"log"

	"github.com/speedianet/sam/src/domain/repository"
	"github.com/speedianet/sam/src/domain/valueObject"
)

func DeleteUser(
	accQueryRepo repository.AccQueryRepo,
	accCmdRepo repository.AccCmdRepo,
	userId valueObject.UserId,
) {
	_, err := accQueryRepo.GetById(userId)
	if err != nil {
		panic("UserIdDoesNotExist")
	}

	err = accCmdRepo.Delete(userId)
	if err != nil {
		panic("UserDeleteError")
	}

	log.Printf("UserId '%v' deleted.", userId)
}
