package useCase

import (
	"log/slog"

	tkDto "github.com/goinfinite/tk/src/domain/dto"
	tkRepository "github.com/goinfinite/tk/src/domain/repository"
)

func ReadActivityRecords(
	activityRecordQueryRepo tkRepository.ActivityRecordQueryRepo,
	readDto tkDto.ReadActivityRecordsRequest,
) (tkDto.ReadActivityRecordsResponse, error) {
	responseDto, err := activityRecordQueryRepo.Read(readDto)
	if err != nil {
		slog.Error("ReadActivityRecordsInfraError", slog.String("err", err.Error()))
		return responseDto, err
	}

	return responseDto, nil
}
