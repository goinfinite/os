package useCase

import (
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/repository"
)

func ReadActivityRecords(
	activityRecordQueryRepo repository.ActivityRecordQueryRepo,
	readDto dto.ReadActivityRecords,
) (activityRecords []entity.ActivityRecord) {
	activityRecords, err := activityRecordQueryRepo.Read(readDto)
	if err != nil {
		slog.Error("ReadActivityRecordsInfraError", slog.Any("err", err))
	}

	return activityRecords
}
