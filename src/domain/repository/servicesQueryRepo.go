package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type ServicesQueryRepo interface {
	Read(
		readDto dto.ReadInstalledServicesItemsRequest,
	) (dto.ReadInstalledServicesItemsResponse, error)
	ReadByName(name valueObject.ServiceName) (entity.InstalledService, error)
	ReadInstallables() ([]entity.InstallableService, error)
}
