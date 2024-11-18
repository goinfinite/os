package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
)

type ServicesQueryRepo interface {
	ReadInstalledItems(
		readDto dto.ReadInstalledServicesItemsRequest,
	) (dto.ReadInstalledServicesItemsResponse, error)
	ReadOneInstalledItem(
		readDto dto.ReadInstalledServicesItemsRequest,
	) (dto.InstalledServiceWithMetrics, error)
	ReadInstallableItems(
		readDto dto.ReadInstallableServicesItemsRequest,
	) (dto.ReadInstallableServicesItemsResponse, error)
}
