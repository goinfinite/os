package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type ServicesQueryRepo interface {
	IsInstalled(valueObject.ServiceName) (bool, error)
	ReadInstalledItems(
		dto.ReadInstalledServicesItemsRequest,
	) (dto.ReadInstalledServicesItemsResponse, error)
	ReadFirstInstalledItem(
		dto.ReadFirstInstalledServiceItemsRequest,
	) (entity.InstalledService, error)
	ReadInstallableItems(
		dto.ReadInstallableServicesItemsRequest,
	) (dto.ReadInstallableServicesItemsResponse, error)
}
