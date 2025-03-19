package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

var FilesDefaultPagination dto.Pagination = dto.Pagination{
	PageNumber:   0,
	ItemsPerPage: 30,
}

func ReadFiles(
	filesQueryRepo repository.FilesQueryRepo,
	requestDto dto.ReadFilesRequest,
) (responseDto dto.ReadFilesResponse, err error) {
	responseDto, err = filesQueryRepo.Read(requestDto)
	if err != nil {
		slog.Error("ReadFilesError", slog.String("err", err.Error()))
		return responseDto, errors.New("ReadFilesInfraError")
	}

	return responseDto, nil
}
