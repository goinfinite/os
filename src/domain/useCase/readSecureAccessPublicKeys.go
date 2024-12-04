package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

var SecureAccessPublicKeysDefaultPagination dto.Pagination = dto.Pagination{
	PageNumber:   0,
	ItemsPerPage: 10,
}

func ReadSecureAccessPublicKeys(
	secureAccessKeyQueryRepo repository.SecureAccessKeyQueryRepo,
	requestDto dto.ReadSecureAccessPublicKeysRequest,
) (secureAccessPublicKeys dto.ReadSecureAccessPublicKeysResponse, err error) {
	secureAccessPublicKeys, err = secureAccessKeyQueryRepo.Read(requestDto)
	if err != nil {
		slog.Error("ReadSecureAccessPublicKeysError", slog.Any("error", err))
		return secureAccessPublicKeys, errors.New(
			"ReadSecureAccessPublicKeysInfraError",
		)
	}

	return secureAccessPublicKeys, nil
}
