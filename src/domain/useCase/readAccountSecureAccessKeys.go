package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
)

func ReadAccountSecureAccessKeys(
	accountQueryRepo repository.AccountQueryRepo,
	accountId valueObject.AccountId,
) (secureAccessKeys []entity.SecureAccessKey, err error) {
	secureAccessKeys, err = accountQueryRepo.ReadSecureAccessKeys(accountId)
	if err != nil {
		slog.Error("ReadAccountSecureAccessKeysInfraError", slog.Any("error", err))
		return secureAccessKeys, errors.New("ReadAccountSecureAccessKeysInfraError")
	}

	return secureAccessKeys, nil
}
