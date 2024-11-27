package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/repository"
)

func ReadCrons(
	cronQueryRepo repository.CronQueryRepo,
) ([]entity.Cron, error) {
	cronsList, err := cronQueryRepo.Read()
	if err != nil {
		slog.Error("ReadCronsError", slog.Any("err", err))
		return cronsList, errors.New("ReadCronsInfraError")
	}

	return cronsList, err
}
