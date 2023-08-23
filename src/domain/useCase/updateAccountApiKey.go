package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/sam/src/domain/dto"
	"github.com/speedianet/sam/src/domain/repository"
	"github.com/speedianet/sam/src/domain/valueObject"
)

func UpdateAccountApiKey(
	accQueryRepo repository.AccQueryRepo,
	accCmdRepo repository.AccCmdRepo,
	updateAccountDto dto.UpdateAccount,
) (valueObject.AccessTokenStr, error) {
	_, err := accQueryRepo.GetById(updateAccountDto.AccountId)
	if err != nil {
		return "", errors.New("AccountNotFound")
	}

	newKey, err := accCmdRepo.UpdateApiKey(updateAccountDto.AccountId)
	if err != nil {
		return "", errors.New("UpdateAccountApiKeyError")
	}

	log.Printf("AccountId '%v' api key updated.", updateAccountDto.AccountId)

	return newKey, nil
}
