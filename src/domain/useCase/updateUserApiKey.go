package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/sam/src/domain/dto"
	"github.com/speedianet/sam/src/domain/repository"
	"github.com/speedianet/sam/src/domain/valueObject"
)

func UpdateUserApiKey(
	accQueryRepo repository.AccQueryRepo,
	accCmdRepo repository.AccCmdRepo,
	updateUserDto dto.UpdateUser,
) (valueObject.AccessTokenStr, error) {
	_, err := accQueryRepo.GetById(updateUserDto.UserId)
	if err != nil {
		return "", errors.New("UserNotFound")
	}

	newKey, err := accCmdRepo.UpdateApiKey(updateUserDto.UserId)
	if err != nil {
		return "", errors.New("UpdateUserApiKeyError")
	}

	log.Printf("UserId '%v' api key updated.", updateUserDto.UserId)

	return newKey, nil
}
