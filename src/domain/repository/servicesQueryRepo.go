package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
)

type ServicesQueryRepo interface {
	ReadInstalledItems(
		readDto dto.ReadInstalledServicesItemsRequest,
	) (dto.ReadInstalledServicesItemsResponse, error)
	ReadUniqueInstalledItem(
		readDto dto.ReadInstalledServicesItemsRequest,
	) (dto.InstalledServiceWithMetrics, error)
	ReadInstallableItems() ([]entity.InstallableService, error)
}
