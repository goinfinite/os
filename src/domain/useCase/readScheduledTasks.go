package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
	tkDto "github.com/goinfinite/tk/src/domain/dto"
)

var ScheduledTasksDefaultPagination tkDto.Pagination = tkDto.Pagination{
	PageNumber:   0,
	ItemsPerPage: 10,
}

func ReadScheduledTasks(
	scheduledTaskQueryRepo repository.ScheduledTaskQueryRepo,
	readDto dto.ReadScheduledTasksRequest,
) (responseDto dto.ReadScheduledTasksResponse, err error) {
	responseDto, err = scheduledTaskQueryRepo.Read(readDto)
	if err != nil {
		slog.Error("ReadTasksInfraError", slog.String("err", err.Error()))
		return responseDto, errors.New("ReadTasksInfraError")
	}

	return responseDto, nil
}
