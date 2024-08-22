package useCase

import (
	"log/slog"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/repository"
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
