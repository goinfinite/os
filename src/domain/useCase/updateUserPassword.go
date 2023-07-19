package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/sam/src/domain/dto"
	"github.com/speedianet/sam/src/domain/repository"
)

func UpdateUserPassword(
	accQueryRepo repository.AccQueryRepo,
	accCmdRepo repository.AccCmdRepo,
	updateUserDto dto.UpdateUser,
) error {
	_, err := accQueryRepo.GetById(updateUserDto.UserId)
	if err != nil {
		return errors.New("UserNotFound")
	}

	err = accCmdRepo.UpdatePassword(
		updateUserDto.UserId,
		*updateUserDto.Password,
	)
	if err != nil {
		return errors.New("UpdateUserPasswordError")
	}

	log.Printf("UserId '%v' password updated.", updateUserDto.UserId)

	return nil
}
