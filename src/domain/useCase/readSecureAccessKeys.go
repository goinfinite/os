package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

var SecureAccessKeysDefaultPagination dto.Pagination = dto.Pagination{
	PageNumber:   0,
	ItemsPerPage: 10,
}

func ReadSecureAccessKeys(
	secureAccessKeyQueryRepo repository.SecureAccessKeyQueryRepo,
	requestDto dto.ReadSecureAccessKeysRequest,
) (secureAccessKeys dto.ReadSecureAccessKeysResponse, err error) {
	secureAccessKeys, err = secureAccessKeyQueryRepo.Read(requestDto)
	if err != nil {
		slog.Error("ReadSecureAccessKeysError", slog.Any("error", err))
		return secureAccessKeys, errors.New("ReadSecureAccessKeysInfraError")
	}

	return secureAccessKeys, nil
}
