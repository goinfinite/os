package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

func UpdateCron(
	cronQueryRepo repository.CronQueryRepo,
	cronCmdRepo repository.CronCmdRepo,
	updateCron dto.UpdateCron,
) error {
	_, err := cronQueryRepo.ReadById(updateCron.Id)
	if err != nil {
		slog.Error("CronNotFound", slog.Any("err", err))
		return errors.New("CronNotFound")
	}

	err = cronCmdRepo.Update(updateCron)
	if err != nil {
		slog.Error("UpdateCronError", slog.Any("err", err))
		return errors.New("UpdateCronInfraError")
	}

	slog.Info("CronUpdated", slog.Uint64("cronId", updateCron.Id.Uint64()))

	return nil
}
