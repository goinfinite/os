package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

var ScheduledTasksDefaultPagination dto.Pagination = dto.Pagination{
	PageNumber:   0,
	ItemsPerPage: 10,
}

func ReadScheduledTasks(
	scheduledTaskQueryRepo repository.ScheduledTaskQueryRepo,
	readDto dto.ReadScheduledTasksRequest,
) (responseDto dto.ReadScheduledTasksResponse, err error) {
	responseDto, err = scheduledTaskQueryRepo.Read(readDto)
	if err != nil {
		slog.Error("ReadTasksInfraError", slog.Any("error", err))
		return responseDto, errors.New("ReadTasksInfraError")
	}

	return responseDto, nil
}
