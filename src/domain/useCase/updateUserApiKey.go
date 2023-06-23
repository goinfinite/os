package useCase

import (
	"log"

	"github.com/speedianet/sam/src/domain/dto"
	"github.com/speedianet/sam/src/domain/repository"
	"github.com/speedianet/sam/src/domain/valueObject"
)

func UpdateUserApiKey(
	accQueryRepo repository.AccQueryRepo,
	accCmdRepo repository.AccCmdRepo,
	updateUserDto dto.UpdateUser,
) valueObject.AccessTokenStr {
	_, err := accQueryRepo.GetById(updateUserDto.UserId)
	if err != nil {
		panic("UserIdDoesNotExist")
	}

	newKey, err := accCmdRepo.UpdateApiKey(updateUserDto.UserId)
	if err != nil {
		panic("UserApiKeyUpdateError")
	}

	log.Printf("UserId '%v' api key updated.", updateUserDto.UserId)
	return newKey
}
