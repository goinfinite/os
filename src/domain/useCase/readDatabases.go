package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

var DatabasesDefaultPagination dto.Pagination = dto.Pagination{
	PageNumber:   0,
	ItemsPerPage: 10,
}

func ReadDatabases(
	databaseQueryRepo repository.DatabaseQueryRepo,
	requestDto dto.ReadDatabasesRequest,
) (responseDto dto.ReadDatabasesResponse, err error) {
	responseDto, err = databaseQueryRepo.Read(requestDto)
	if err != nil {
		slog.Error("ReadDatabasesError", slog.String("err", err.Error()))
		return responseDto, errors.New("ReadDatabasesInfraError")
	}

	return responseDto, nil
}
