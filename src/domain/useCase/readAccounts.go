package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
	tkDto "github.com/goinfinite/tk/src/domain/dto"
)

var AccountsDefaultPagination tkDto.Pagination = tkDto.Pagination{
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
