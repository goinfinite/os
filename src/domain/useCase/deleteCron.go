package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

func DeleteCron(
	cronQueryRepo repository.CronQueryRepo,
	cronCmdRepo repository.CronCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	deleteDto dto.DeleteCron,
) error {
	if deleteDto.Id == nil {
		if deleteDto.Comment == nil {
			return errors.New("CronIdOrCommentRequired")
		}

		readFirstRequestDto := dto.ReadCronsRequest{
			CronComment: deleteDto.Comment,
		}
		cron, err := cronQueryRepo.ReadFirst(readFirstRequestDto)
		if err != nil {
			slog.Error("ReadCronToDeleteError", slog.Any("err", err))
			return errors.New("ReadCronToDeleteInfraError")
		}
		deleteDto.Id = &cron.Id
	}

	err := cronCmdRepo.Delete(*deleteDto.Id)
	if err != nil {
		slog.Error("DeleteCronError", slog.Any("err", err))
		return errors.New("DeleteCronInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).DeleteCron(deleteDto)

	return nil
}
