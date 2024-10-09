package dto

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type InstalledServiceWithMetrics struct {
	entity.InstalledService
	Metrics *valueObject.ServiceMetrics `json:"metrics"`
}

func NewInstalledServiceWithMetrics(
	installedService entity.InstalledService,
	metrics *valueObject.ServiceMetrics,
) InstalledServiceWithMetrics {
	return InstalledServiceWithMetrics{
		InstalledService: installedService,
		Metrics:          metrics,
	}
}
