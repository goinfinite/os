package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
)

func ReadSecureAccessKeys(
	secureAccessKeyQueryRepo repository.SecureAccessKeyQueryRepo,
	accountId valueObject.AccountId,
) (secureAccessKeys []entity.SecureAccessKey, err error) {
	secureAccessKeys, err = secureAccessKeyQueryRepo.Read(accountId)
	if err != nil {
		slog.Error("ReadSecureAccessKeysInfraError", slog.Any("error", err))
		return secureAccessKeys, errors.New("ReadSecureAccessKeysInfraError")
	}

	return secureAccessKeys, nil
}
