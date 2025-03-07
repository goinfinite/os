package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

var CronsDefaultPagination dto.Pagination = dto.Pagination{
	PageNumber:   0,
	ItemsPerPage: 10,
}

func ReadCrons(
	cronQueryRepo repository.CronQueryRepo,
	requestDto dto.ReadCronsRequest,
) (responseDto dto.ReadCronsResponse, err error) {
	responseDto, err = cronQueryRepo.Read(requestDto)
	if err != nil {
		slog.Error("ReadCronsError", slog.String("err", err.Error()))
		return responseDto, errors.New("ReadCronsInfraError")
	}

	return responseDto, err
}
