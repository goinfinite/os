package useCase

import (
	"errors"
	"log"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
)

func AccountLookup(
	accQueryRepo repository.AccQueryRepo,
	accountId *valueObject.AccountId,
	username *valueObject.Username,
) (accountEntity entity.Account, err error) {
	if accountId == nil && username == nil {
		return accountEntity, errors.New("AccountIdOrUsernameRequired")
	}

	if accountId == nil {
		accountEntity, err := accQueryRepo.GetByUsername(*username)
		if err != nil {
			return accountEntity, errors.New("AccountNotFound")
		}
		return accountEntity, nil
	}

	accountEntity, err = accQueryRepo.GetById(*accountId)
	if err != nil {
		return accountEntity, errors.New("AccountNotFound")
	}

	return accountEntity, nil
}

func UpdateAccountApiKey(
	accQueryRepo repository.AccQueryRepo,
	accCmdRepo repository.AccCmdRepo,
	updateDto dto.UpdateAccount,
) (accessToken valueObject.AccessTokenStr, err error) {
	accountEntity, err := AccountLookup(accQueryRepo, updateDto.Id, updateDto.Username)
	if err != nil {
		return accessToken, err
	}

	newKey, err := accCmdRepo.UpdateApiKey(accountEntity.Id)
	if err != nil {
		log.Printf("UpdateAccountApiKeyError: %v", err)
		return accessToken, errors.New("UpdateAccountApiKeyInfraError")
	}

	log.Printf("AccountId '%v' api key updated.", accountEntity.Id)
	return newKey, nil
}

func UpdateAccountPassword(
	accQueryRepo repository.AccQueryRepo,
	accCmdRepo repository.AccCmdRepo,
	updateDto dto.UpdateAccount,
) error {
	accountEntity, err := AccountLookup(accQueryRepo, updateDto.Id, updateDto.Username)
	if err != nil {
		return err
	}

	err = accCmdRepo.UpdatePassword(accountEntity.Id, *updateDto.Password)
	if err != nil {
		log.Printf("UpdateAccountPasswordError: %v", err)
		return errors.New("UpdateAccountPasswordInfraError")
	}

	log.Printf("AccountId '%v' password updated.", accountEntity.Id)
	return nil
}
