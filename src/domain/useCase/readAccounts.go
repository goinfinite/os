package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

func ReadAccounts(
	accountQueryRepo repository.AccountQueryRepo,
	requestDto dto.ReadAccountsRequest,
) (responseDto dto.ReadAccountsResponse, err error) {
	responseDto, err = accountQueryRepo.Read(requestDto)
	if err != nil {
		slog.Error("ReadAccountsInfraError", slog.Any("error", err))
		return responseDto, errors.New("ReadAccountsInfraError")
	}

	return responseDto, nil
}
