package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
)

func UpdateAccountApiKey(
	accQueryRepo repository.AccQueryRepo,
	accCmdRepo repository.AccCmdRepo,
	updateAccountDto dto.UpdateAccount,
) (valueObject.AccessTokenStr, error) {
	_, err := accQueryRepo.GetById(updateAccountDto.Id)
	if err != nil {
		return "", errors.New("AccountNotFound")
	}

	newKey, err := accCmdRepo.UpdateApiKey(updateAccountDto.Id)
	if err != nil {
		return "", errors.New("UpdateAccountApiKeyError")
	}

	log.Printf("AccountId '%v' api key updated.", updateAccountDto.Id)

	return newKey, nil
}
