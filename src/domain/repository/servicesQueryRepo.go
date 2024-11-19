package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
)

type ServicesQueryRepo interface {
	ReadInstalledItems(
		requestDto dto.ReadInstalledServicesItemsRequest,
	) (dto.ReadInstalledServicesItemsResponse, error)
	ReadOneInstalledItem(
		requestDto dto.ReadInstalledServicesItemsRequest,
	) (dto.InstalledServiceWithMetrics, error)
	ReadInstallableItems(
		requestDto dto.ReadInstallableServicesItemsRequest,
	) (dto.ReadInstallableServicesItemsResponse, error)
}
