package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
)

func ReadSecureAccessKeys(
	accountQueryRepo repository.AccountQueryRepo,
	secureAccessKeyQueryRepo repository.SecureAccessKeyQueryRepo,
	accountId valueObject.AccountId,
) (secureAccessKeys []entity.SecureAccessKey, err error) {
	_, err = accountQueryRepo.ReadById(accountId)
	if err != nil {
		return secureAccessKeys, errors.New("AccountNotFound")
	}

	secureAccessKeys, err = secureAccessKeyQueryRepo.Read(accountId)
	if err != nil {
		slog.Error("ReadSecureAccessKeysInfraError", slog.Any("error", err))
		return secureAccessKeys, errors.New("ReadSecureAccessKeysInfraError")
	}

	return secureAccessKeys, nil
}
