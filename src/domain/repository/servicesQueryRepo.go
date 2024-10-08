package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type ServicesQueryRepo interface {
	Read() ([]entity.InstalledService, error)
	ReadByName(name valueObject.ServiceName) (entity.InstalledService, error)
	ReadWithMetrics() ([]dto.InstalledServiceWithMetrics, error)
	ReadInstallables() ([]entity.InstallableService, error)
}
