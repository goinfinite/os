package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
)

type ServicesQueryRepo interface {
	ReadInstalledItems(
		dto.ReadInstalledServicesItemsRequest,
	) (dto.ReadInstalledServicesItemsResponse, error)
	ReadFirstInstalledItem(
		dto.ReadInstalledServicesItemsRequest,
	) (entity.InstalledService, error)
	ReadInstallableItems(
		dto.ReadInstallableServicesItemsRequest,
	) (dto.ReadInstallableServicesItemsResponse, error)
}
