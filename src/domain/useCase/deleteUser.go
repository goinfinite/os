package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/sam/src/domain/repository"
	"github.com/speedianet/sam/src/domain/valueObject"
)

func DeleteUser(
	accQueryRepo repository.AccQueryRepo,
	accCmdRepo repository.AccCmdRepo,
	userId valueObject.UserId,
) error {
	_, err := accQueryRepo.GetById(userId)
	if err != nil {
		return errors.New("UserNotFound")
	}

	err = accCmdRepo.Delete(userId)
	if err != nil {
		return errors.New("DeleteUserError")
	}

	log.Printf("UserId '%v' deleted.", userId)

	return nil
}
