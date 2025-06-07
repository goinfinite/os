package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

var AccountsDefaultPagination dto.Pagination = dto.Pagination{
	PageNumber:   0,
	ItemsPerPage: 10,
}

func ReadAccounts(
	accountQueryRepo repository.AccountQueryRepo,
	requestDto dto.ReadAccountsRequest,
) (responseDto dto.ReadAccountsResponse, err error) {
	responseDto, err = accountQueryRepo.Read(requestDto)
	if err != nil {
		slog.Error("ReadAccountsInfraError", slog.String("err", err.Error()))
		return responseDto, errors.New("ReadAccountsInfraError")
	}

	return responseDto, nil
}
