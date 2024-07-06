package repository

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type ServicesQueryRepo interface {
	Read() ([]entity.InstalledService, error)
	ReadWithMetrics() ([]dto.InstalledServiceWithMetrics, error)
	ReadByName(name valueObject.ServiceName) (entity.InstalledService, error)
	ReadInstallables() ([]entity.InstallableService, error)
}
