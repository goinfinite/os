package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

var SslPairsDefaultPagination dto.Pagination = dto.Pagination{
	PageNumber:   0,
	ItemsPerPage: 10,
}

func ReadSslPairs(
	sslQueryRepo repository.SslQueryRepo,
	requestDto dto.ReadSslPairsRequest,
) (responseDto dto.ReadSslPairsResponse, err error) {
	responseDto, err = sslQueryRepo.Read(requestDto)
	if err != nil {
		slog.Error("ReadSslPairsError", slog.String("err", err.Error()))
		return responseDto, errors.New("ReadSslPairsInfraError")
	}

	return responseDto, nil
}
