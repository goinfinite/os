package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/repository"
)

func ReadO11yOverview(
	o11yQueryRepo repository.O11yQueryRepo,
	withResourceUsage bool,
) (entity.O11yOverview, error) {
	o11yOverviewEntity, err := o11yQueryRepo.ReadOverview(withResourceUsage)
	if err != nil {
		slog.Error("ReadOverviewError", slog.String("err", err.Error()))
		return o11yOverviewEntity, errors.New("ReadOverviewInfraError")
	}

	return o11yOverviewEntity, nil
}
