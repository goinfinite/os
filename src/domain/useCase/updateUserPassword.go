package useCase

import (
	"log"

	"github.com/speedianet/sam/src/domain/dto"
	"github.com/speedianet/sam/src/domain/repository"
)

func UpdateUserPassword(
	accQueryRepo repository.AccQueryRepo,
	accCmdRepo repository.AccCmdRepo,
	updateUserDto dto.UpdateUser,
) {
	_, err := accQueryRepo.GetById(updateUserDto.UserId)
	if err != nil {
		panic("UserIdDoesNotExist")
	}

	err = accCmdRepo.UpdatePassword(
		updateUserDto.UserId,
		*updateUserDto.Password,
	)
	if err != nil {
		panic("PasswordUpdateError")
	}

	log.Printf("UserId '%v' password updated.", updateUserDto.UserId)
}
